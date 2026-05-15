import { supabase } from './supabase'

const appRedirectUrl = () => `${window.location.origin}${window.location.pathname}`

export const signUpWithEmail = async (email, password) =>
  supabase.auth.signUp({ email, password, options: { emailRedirectTo: appRedirectUrl() } })

export const signInWithEmail = async (email, password) =>
  supabase.auth.signInWithPassword({ email, password })

export const signUpWithPhone = async (phone) =>
  supabase.auth.signInWithOtp({ phone })

export const verifyPhoneOtp = async (phone, token) =>
  supabase.auth.verifyOtp({ phone, token, type: 'sms' })

export const signInWithGoogle = async () =>
  supabase.auth.signInWithOAuth({ provider: 'google', options: { redirectTo: appRedirectUrl() } })

export const signOut = async () => supabase.auth.signOut()

export const getSession = async () => (await supabase.auth.getSession()).data.session
