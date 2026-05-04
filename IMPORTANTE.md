# =============================================================================
  Informacion sobre modelos LLM por apikey 
# =============================================================================
Usar OpenRouter con modelos free. Verificados Mayo 2026:

| Uso | Modelo | Benchmark Score | Contexto |
|-----|--------|----------------|----------|
| **PRINCIPAL (tool calling)** | `google/gemma-4-31b-it:free` | GPQA 85.7%, MMLU-Pro 85.2% | 256K |
| VISION/Multimodal | `nvidia/nemotron-nano-12b-v2-vl:free` | — | — |
| RESPUESTAS RAPIDAS | `inclusionai/ling-2.6-1t:free` | — | — |

Gemma-4-31B-it es el modelo recomendado para ION:
- Superior en agente/razonamiento vs Nemotron-3-Super-120B
- Multimodal nativo (texto + imagen + video)
- Tool calling nativo para flujos agenticos
- Apache 2.0 license
- Benchmarks: GPQA 85.7%, MMLU-Pro 85.2%, Terminal-Bench Hard 36.4%

verificar siempre que los modelos sean :free y que sigan siendo los mejores actuales
confirmar con internet, prohibido asumir o alucinar modelos

# =============================================================================
  Informacion sobre el servidor cx23 hetzner de produccion
# ============================================================================= 
iodesk-server
IP= 178.104.36.144

# =============================================================================
  Informacion correos y agentes
# ============================================================================= 
ion@iodevs.net es un correo que funciona gracias al reenvio de cloudflare, el correo es el.agente.ion@gmail.com
ION es el interlocutor entre usuarios humanos y los agentes en localAGI
ION solo debe responder a correos que llegan desde los siguientes correos:
    - @ionet.cl
    - @iodevs.net
    - el.agente.ion@gmail.com
    - ionet.ventas@gmail.com
