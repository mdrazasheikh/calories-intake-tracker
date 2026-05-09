import { serve } from 'https://deno.land/std@0.224.0/http/server.ts'

const corsHeaders = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Headers': 'authorization, x-client-info, apikey, content-type'
}

serve(async (req) => {
  if (req.method === 'OPTIONS') return new Response('ok', { headers: corsHeaders })
  if (req.method !== 'POST') return Response.json({ error: 'method not allowed' }, { status: 405, headers: corsHeaders })

  const { q } = await req.json()
  if (!q || typeof q !== 'string') return Response.json({ error: 'q is required' }, { status: 400, headers: corsHeaders })

  const endpoint = `https://world.openfoodfacts.org/api/v2/foods/search?search_terms=${encodeURIComponent(q)}&fields=product_name,brands,code,nutriments&json=true`
  const res = await fetch(endpoint)
  if (!res.ok) return Response.json({ error: 'upstream error' }, { status: 502, headers: corsHeaders })

  const json = await res.json()
  const items = (json.products ?? [])
    .filter((p: { product_name?: string }) => p.product_name)
    .slice(0, 20)
    .map((p: any) => ({
      name: p.product_name,
      brand: p.brands,
      barcode: p.code,
      calories: p.nutriments?.['energy-kcal_100g'] ?? 0,
      serving: '100g',
      source: 'openfoodfacts'
    }))

  return Response.json(items, { headers: corsHeaders })
})
