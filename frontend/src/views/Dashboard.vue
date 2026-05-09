<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../services/api'
import { getSession, signOut } from '../services/auth'

const query = ref(''); const barcode = ref(''); const foods = ref([]); const recommendation = ref(null); const session = ref(null); const error = ref('')
const profile = ref({ sex:'female', age:30, height_cm:165, weight_kg:65, activity_level:'moderate', goal_type:'lose' })
onMounted(async()=>{ session.value = await getSession() })
async function searchFood(){ try { foods.value = await api.searchFoods(query.value); error.value='' } catch(e){ error.value=e.message } }
async function scanBarcode(){ try { foods.value=[await api.lookupBarcode(barcode.value)]; error.value='' } catch(e){ error.value=e.message } }
async function recommend(){ try { recommendation.value=await api.recommend(profile.value); error.value='' } catch(e){ error.value=e.message } }
async function logout(){ await signOut(); session.value = null }
</script>
<template>
  <main class="container">
    <div class="row"><h1>Calorie Tracker</h1><button v-if="session" @click="logout">Logout</button></div>
    <p v-if="!session">You are not logged in. Go to <router-link to="/auth">Auth</router-link>.</p>
    <p v-if="error" class="error">{{ error }}</p>
    <section><h2>Food Lookup</h2><input v-model="query" placeholder="Search food"/><button @click="searchFood">Search</button>
    <input v-model="barcode" placeholder="Barcode"/><button @click="scanBarcode">Lookup Barcode</button></section>
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
.container{max-width:720px;margin:auto;padding:1rem;font-family:system-ui} input,button,select{padding:.6rem;margin:.2rem;width:100%} .grid{display:grid;grid-template-columns:1fr 1fr 1fr;gap:.3rem}.row{display:flex;align-items:center;justify-content:space-between}.error{color:#b00020} @media(max-width:640px){.grid{grid-template-columns:1fr}}
</style>
