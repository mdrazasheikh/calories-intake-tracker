import { serve } from 'https://deno.land/std@0.224.0/http/server.ts'

serve(async (req) => {
  const url = new URL(req.url)
  const barcode = url.searchParams.get('barcode') ?? ''
  const endpoint = `https://world.openfoodfacts.org/api/v2/product/${encodeURIComponent(barcode)}?fields=product_name,brands,code,nutriments`
  const res = await fetch(endpoint)
  const json = await res.json()
  const p = json.product
  if (!p?.product_name) return Response.json({ error: 'barcode not found' }, { status: 404 })

  return Response.json({
    name: p.product_name,
    brand: p.brands,
    barcode: p.code,
    calories: p.nutriments?.['energy-kcal_100g'] ?? 0,
    serving: '100g',
    source: 'openfoodfacts'
  })
})
