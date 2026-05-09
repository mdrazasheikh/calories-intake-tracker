import { serve } from 'https://deno.land/std@0.224.0/http/server.ts'

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type'
}

serve(async (req) => {
  if (req.method === 'OPTIONS') return new Response('ok', { headers: corsHeaders })
  if (req.method !== 'POST') return Response.json({ error: 'method not allowed' }, { status: 405, headers: corsHeaders })

  const { barcode } = await req.json()
  if (!barcode || typeof barcode !== 'string') return Response.json({ error: 'barcode is required' }, { status: 400, headers: corsHeaders })

  const endpoint = `https://world.openfoodfacts.org/api/v2/product/${encodeURIComponent(barcode)}?fields=product_name,brands,code,nutriments`
  const res = await fetch(endpoint)
  if (!res.ok) return Response.json({ error: 'upstream error' }, { status: 502, headers: corsHeaders })

  const json = await res.json()
  const p = json.product
  if (!p?.product_name) return Response.json({ error: 'barcode not found' }, { status: 404, headers: corsHeaders })

  return Response.json({
    name: p.product_name,
    brand: p.brands,
    barcode: p.code,
    calories: p.nutriments?.['energy-kcal_100g'] ?? 0,
    serving: '100g',
    source: 'openfoodfacts'
  }, { headers: corsHeaders })
})
