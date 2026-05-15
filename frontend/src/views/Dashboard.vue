<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { api } from '../services/api'
import { getSession, signOut } from '../services/auth'

const query = ref(''); const barcode = ref(''); const foods = ref([]); const recommendation = ref(null); const session = ref(null); const error = ref('')
const profile = ref({ sex:'female', age:30, height_cm:165, weight_kg:65, activity_level:'moderate', goal_type:'lose' })
const scannerVideo = ref(null); const captureCanvas = ref(null); const scannerOpen = ref(false); const scannerActive = ref(false); const scannerMessage = ref('')
const decodedCode = ref(''); const decodedFormat = ref(''); const decodingFrame = ref(false)
const lookingUpBarcode = ref(false)
let barcodeReader = null
let cameraStream = null
const barcodeFormatLabels = {
  2: 'Code 39',
  4: 'Code 128',
  6: 'EAN-8',
  7: 'EAN-13',
  8: 'ITF',
  11: 'QR Code',
  14: 'UPC-A',
  15: 'UPC-E'
}
onMounted(async()=>{ session.value = await getSession() })
async function searchFood(){ try { foods.value = await api.searchFoods(query.value); error.value='' } catch(e){ error.value=e.message } }
async function lookupBarcode(value){
  const code = value.trim()
  if (!code) {
    error.value = 'Enter or scan a barcode.'
    return
  }

  lookingUpBarcode.value = true
  try {
    foods.value = [await api.lookupBarcode(code)]
    barcode.value = code
    error.value = ''
  } catch(e) {
    error.value = e.message
  } finally {
    lookingUpBarcode.value = false
  }
}
async function scanBarcode(){ await lookupBarcode(barcode.value) }
async function recommend(){ try { recommendation.value=await api.recommend(profile.value); error.value='' } catch(e){ error.value=e.message } }
async function logout(){ await signOut(); session.value = null }
function stopCameraScanner(){
  if (cameraStream) {
    cameraStream.getTracks().forEach((track) => track.stop())
    cameraStream = null
  }

  if (scannerVideo.value) {
    scannerVideo.value.srcObject = null
  }

  scannerActive.value = false
  scannerOpen.value = false
  scannerMessage.value = ''
}
async function getBarcodeReader(){
  if (barcodeReader) return barcodeReader

  const [{ BrowserMultiFormatReader }, { BarcodeFormat, DecodeHintType }] = await Promise.all([
    import('@zxing/browser'),
    import('@zxing/library')
  ])
  const hints = new Map()
  hints.set(DecodeHintType.POSSIBLE_FORMATS, [
    BarcodeFormat.EAN_13,
    BarcodeFormat.EAN_8,
    BarcodeFormat.UPC_A,
    BarcodeFormat.UPC_E,
    BarcodeFormat.CODE_128,
    BarcodeFormat.CODE_39,
    BarcodeFormat.ITF,
    BarcodeFormat.QR_CODE
  ])
  hints.set(DecodeHintType.TRY_HARDER, true)
  barcodeReader = new BrowserMultiFormatReader(hints)
  return barcodeReader
}
async function startCameraScanner(){
  if (!navigator.mediaDevices?.getUserMedia) {
    error.value = 'Camera scanning requires a browser with camera support.'
    return
  }

  stopCameraScanner()
  scannerOpen.value = true
  decodedCode.value = ''
  decodedFormat.value = ''
  scannerMessage.value = 'Frame the barcode or QR code, then tap Capture & Decode.'
  error.value = ''
  await nextTick()

  try {
    const constraints = {
      video: {
        facingMode: { ideal: 'environment' },
        width: { ideal: 1920 },
        height: { ideal: 1080 }
      },
      audio: false
    }
    cameraStream = await navigator.mediaDevices.getUserMedia(constraints)
    scannerVideo.value.srcObject = cameraStream
    await scannerVideo.value.play()
    scannerActive.value = true
  } catch(e) {
    stopCameraScanner()
    error.value = e?.message || 'Could not start the camera scanner.'
  }
}
async function captureAndDecode(){
  if (!scannerVideo.value || !captureCanvas.value) return

  const video = scannerVideo.value
  if (!video.videoWidth || !video.videoHeight) {
    scannerMessage.value = 'Camera is still warming up. Try again in a moment.'
    return
  }

  decodingFrame.value = true
  error.value = ''
  try {
    const canvas = captureCanvas.value
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    canvas.getContext('2d').drawImage(video, 0, 0, canvas.width, canvas.height)

    const reader = await getBarcodeReader()
    const result = reader.decodeFromCanvas(canvas)
    const text = result.getText()
    const format = result.getBarcodeFormat?.()
    decodedCode.value = text
    decodedFormat.value = barcodeFormatLabels[format] || 'Detected code'
    barcode.value = text
    scannerMessage.value = 'Decoded from captured image.'
  } catch(e) {
    decodedCode.value = ''
    decodedFormat.value = ''
    scannerMessage.value = 'No barcode or QR code found in that image. Hold it steady, fill more of the frame, and capture again.'
  } finally {
    decodingFrame.value = false
  }
}
async function lookupDecodedCode(){ await lookupBarcode(decodedCode.value) }
onUnmounted(stopCameraScanner)
</script>
<template>
  <main class="container">
    <div class="row"><h1>Calorie Tracker</h1><button v-if="session" @click="logout">Logout</button></div>
    <p v-if="!session">You are not logged in. Go to <router-link to="/auth">Auth</router-link>.</p>
    <p v-if="error" class="error">{{ error }}</p>
    <section><h2>Food Lookup</h2><input v-model="query" placeholder="Search food"/><button @click="searchFood">Search</button>
    <div class="barcode-actions"><input v-model="barcode" inputmode="numeric" placeholder="Barcode"/><button @click="scanBarcode" :disabled="lookingUpBarcode">Lookup Barcode</button><button type="button" @click="startCameraScanner">Scan with Camera</button></div>
    <div v-if="scannerOpen" class="scanner">
      <video ref="scannerVideo" muted playsinline></video>
      <canvas ref="captureCanvas" aria-hidden="true"></canvas>
      <p>{{ scannerMessage }}</p>
      <div v-if="decodedCode" class="decoded-result">
        <strong>{{ decodedFormat }}</strong>
        <span>{{ decodedCode }}</span>
        <button type="button" @click="lookupDecodedCode" :disabled="lookingUpBarcode">Lookup Decoded Code</button>
      </div>
      <button type="button" @click="captureAndDecode" :disabled="decodingFrame || !scannerActive">{{ decodingFrame ? 'Decoding...' : 'Capture & Decode' }}</button>
      <button type="button" @click="stopCameraScanner">{{ scannerActive ? 'Stop Camera' : 'Close Scanner' }}</button>
    </div>
    </section>
    <ul><li v-for="f in foods" :key="f.barcode || f.name">{{ f.name }} - {{ Math.round(f.calories) }} kcal / {{f.serving}}</li></ul>
    <section><h2>Calorie Recommendation</h2><div class="grid"><input v-model.number="profile.age" type="number" placeholder="Age"/><input v-model.number="profile.height_cm" type="number" placeholder="Height cm"/><input v-model.number="profile.weight_kg" type="number" placeholder="Weight kg"/></div>
      <select v-model="profile.sex"><option>female</option><option>male</option></select>
      <select v-model="profile.activity_level"><option value="sedentary">Sedentary</option><option value="light">Light</option><option value="moderate">Moderate</option><option value="active">Active</option></select>
      <select v-model="profile.goal_type"><option value="lose">Lose</option><option value="maintain">Maintain</option><option value="gain">Gain</option></select>
      <button @click="recommend">Recommend</button>
      <p v-if="recommendation">Recommended: {{ Math.round(recommendation.recommended_calories) }} kcal/day</p>
    </section>
  </main>
</template>
<style scoped>
.container{max-width:720px;margin:auto;padding:1rem;font-family:system-ui} input,button,select{padding:.6rem;margin:.2rem;width:100%} button:disabled{opacity:.65;cursor:not-allowed}.grid{display:grid;grid-template-columns:1fr 1fr 1fr;gap:.3rem}.row{display:flex;align-items:center;justify-content:space-between}.barcode-actions{display:grid;grid-template-columns:1fr 1fr 1fr;gap:.3rem}.scanner{margin:.5rem .2rem 0}.scanner video{display:block;width:100%;aspect-ratio:16/9;background:#111;border-radius:8px;object-fit:cover}.scanner canvas{display:none}.scanner p{margin:.5rem 0;color:#333}.decoded-result{display:grid;gap:.35rem;margin:.5rem .2rem;padding:.7rem;border:1px solid #ddd;border-radius:8px;background:#fafafa}.decoded-result span{overflow-wrap:anywhere}.error{color:#b00020} @media(max-width:640px){.grid,.barcode-actions{grid-template-columns:1fr}}
</style>
