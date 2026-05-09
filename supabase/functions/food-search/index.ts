import { serve } from 'https://deno.land/std@0.224.0/http/server.ts'

serve(async (req) => {
  const url = new URL(req.url)
  const q = url.searchParams.get('q') ?? ''
  const endpoint = `https://world.openfoodfacts.org/api/v2/foods/search?search_terms=${encodeURIComponent(q)}&fields=product_name,brands,code,nutriments&json=true`
  const res = await fetch(endpoint)
  const json = await res.json()
  const items = (json.products ?? []).filter((p:any) => p.product_name).slice(0, 20).map((p:any) => ({
    name: p.product_name,
    brand: p.brands,
    barcode: p.code,
    calories: p.nutriments?.['energy-kcal_100g'] ?? 0,
    serving: '100g',
    source: 'openfoodfacts'
  }))
  return Response.json(items)
})
