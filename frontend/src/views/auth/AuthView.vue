<script setup>
import { ref } from 'vue'
import { signInWithEmail, signUpWithEmail, signUpWithPhone, verifyPhoneOtp, signInWithGoogle } from '../../services/auth'

const email = ref(''); const password = ref(''); const phone = ref(''); const otp = ref(''); const message = ref('')

const doEmailSignUp = async () => { const { error } = await signUpWithEmail(email.value, password.value); message.value = error ? error.message : 'Check your email to verify your account.' }
const doEmailSignIn = async () => { const { error } = await signInWithEmail(email.value, password.value); message.value = error ? error.message : 'Signed in.' }
const doPhoneStart = async () => { const { error } = await signUpWithPhone(phone.value); message.value = error ? error.message : 'OTP sent by SMS.' }
const doPhoneVerify = async () => { const { error } = await verifyPhoneOtp(phone.value, otp.value); message.value = error ? error.message : 'Phone verified and signed in.' }
const doGoogle = async () => { const { error } = await signInWithGoogle(); if (error) message.value = error.message }
</script>

<template>
<main class="container">
  <h1>Sign In / Sign Up</h1>
  <input v-model="email" placeholder="Email" type="email" />
  <input v-model="password" placeholder="Password" type="password" />
  <button @click="doEmailSignUp">Sign up with Email</button>
  <button @click="doEmailSignIn">Sign in with Email</button>
  <hr />
  <input v-model="phone" placeholder="Phone (+15551234567)" />
  <button @click="doPhoneStart">Send Phone OTP</button>
  <input v-model="otp" placeholder="OTP code" />
  <button @click="doPhoneVerify">Verify Phone OTP</button>
  <hr />
  <button @click="doGoogle">Continue with Google</button>
  <p>{{ message }}</p>
</main>
</template>

<style scoped>.container{max-width:720px;margin:auto;padding:1rem} input,button{width:100%;padding:.6rem;margin:.25rem 0}</style>
