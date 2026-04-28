"""Test: reproduce cogito pick_tool hang with nemotron streaming + forced tool."""
import os, json, sys, time
from openai import OpenAI

api_key = os.environ.get("OPENAI_API_KEY", "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")
base_url = os.environ.get("OPENAI_BASE_URL", "https://openrouter.ai/api/v1")
model = os.environ.get("MODEL_NAME", "nvidia/nemotron-3-super-120b-a12b:free")

client = OpenAI(api_key=api_key, base_url=base_url)

# The exact pick_tool definition cogito uses
pick_tool = {
    "type": "function",
    "function": {
        "name": "pick_tool",
        "description": "Pick the most appropriate tool to use based on the reasoning.",
        "strict": False,
        "parameters": {
            "type": "object",
            "properties": {
                "tool": {
                    "type": "string",
                    "description": "The tool to use",
                    "enum": ["call_agents"]
                },
                "reasoning": {
                    "type": "string",
                    "description": "The reasoning for the tool choice"
                }
            },
            "required": ["tool"]
        }
    }
}

messages = [
    {"role": "system", "content": "You are an email assistant. Answer emails concisely."},
    {"role": "user", "content": "Subject: Test email\nFrom: user@test.com\n\nThis is a test email asking about password reset."},
    {"role": "assistant", "content": "The user is asking about password reset procedures. I need to determine which tool to use to respond to this inquiry."}
]

print(f"=== Test 2: STREAMING with forced tool (cogito path) ===")
start = time.time()
try:
    stream = client.chat.completions.create(
        model=model,
        messages=messages,
        tools=[pick_tool],
        tool_choice={"type": "function", "function": {"name": "pick_tool"}},
        stream=True,
        stream_options={"include_usage": True},
        timeout=60
    )
    elapsed = time.time() - start
    print(f"Stream opened in {elapsed:.1f}s")

    tool_calls = {}
    last_finish = ""
    chunk_count = 0
    start = time.time()
    for chunk in stream:
        chunk_count += 1
        elapsed = time.time() - start
        if not chunk.choices:
            if chunk.usage:
                print(f"  [{elapsed:.1f}s] usage chunk: {chunk.usage}")
            continue

        delta = chunk.choices[0].delta
        finish = chunk.choices[0].finish_reason

        # Handle reasoning_content safely (may exist in dict form)
        rc = getattr(delta, 'reasoning_content', None) or (delta.model_extra.get('reasoning_content') if hasattr(delta, 'model_extra') and delta.model_extra else None)
        if rc:
            print(f"  [{elapsed:.1f}s] reasoning({len(rc)}c): {rc[:80]}")

        if delta.content:
            print(f"  [{elapsed:.1f}s] content({len(delta.content)}c): {delta.content[:80]}")

        if delta.tool_calls:
            for tc in delta.tool_calls:
                idx = tc.index if tc.index is not None else 0
                if idx not in tool_calls:
                    tool_calls[idx] = {"name": "", "args": ""}
                if tc.function and tc.function.name:
                    tool_calls[idx]["name"] += tc.function.name
                if tc.function and tc.function.arguments:
                    tool_calls[idx]["args"] += tc.function.arguments
                print(f"  [{elapsed:.1f}s] tool[{idx}]: name='{tc.function.name if tc.function else ''}' args='{tc.function.arguments if tc.function else ''}'")

        if finish:
            last_finish = finish
            print(f"  [{elapsed:.1f}s] FINISH: {finish}")

    total = time.time() - start
    print(f"\n=== RESULT ({total:.1f}s, {chunk_count} chunks) ===")
    print(f"Tool calls: {json.dumps(tool_calls, indent=2)}")
    print(f"Finish: {last_finish}")

except Exception as e:
    elapsed = time.time() - start
    print(f"ERROR after {elapsed:.1f}s: {type(e).__name__}: {e}")

print()
# Now a raw-request test to see exact HTTP response
print("=== Test 3: RAW streaming check ===")
import requests
start = time.time()
try:
    resp = requests.post(
        f"{base_url}/chat/completions",
        headers={
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
        },
        json={
            "model": model,
            "messages": messages,
            "tools": [pick_tool],
            "tool_choice": {"type": "function", "function": {"name": "pick_tool"}},
            "stream": True,
        },
        stream=True,
        timeout=60
    )
    elapsed = time.time() - start
    print(f"Stream opened in {elapsed:.1f}s, status={resp.status_code}")

    line_count = 0
    start = time.time()
    for line in resp.iter_lines():
        if line:
            line_count += 1
            elapsed = time.time() - start
            decoded = line.decode('utf-8', errors='replace')
            if line_count <= 20 or 'tool_calls' in decoded.lower() or '[done]' in decoded.lower():
                print(f"  [{elapsed:.1f}s] {decoded[:200]}")

    total = time.time() - start
    print(f"\nTotal: {total:.1f}s, {line_count} lines")

except Exception as e:
    elapsed = time.time() - start
    print(f"ERROR after {elapsed:.1f}s: {type(e).__name__}: {e}")
