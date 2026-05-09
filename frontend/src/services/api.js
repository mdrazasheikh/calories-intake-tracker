import { supabase } from './supabase'

export const api = {
  async searchFoods(q) {
    const { data, error } = await supabase.functions.invoke(`food-search?q=${encodeURIComponent(q)}`, { method: 'GET' })
    if (error) throw error
    return data
  },
  async lookupBarcode(barcode) {
    const { data, error } = await supabase.functions.invoke(`barcode-lookup?barcode=${encodeURIComponent(barcode)}`, { method: 'GET' })
    if (error) throw error
    return data
  },
  async recommend(payload) {
    const { data, error } = await supabase.functions.invoke('calorie-recommendation', { body: payload })
    if (error) throw error
    return data
  }
}
