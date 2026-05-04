const http = require('http');
const https = require('https');
const url = require('url');

const PORT = process.env.PORT || 3001;
const DEEPSEEK_HOST = process.env.DEEPSEEK_HOST || 'api.deepseek.com';
const DEEPSEEK_PROTO = process.env.DEEPSEEK_PROTO || 'https';

const server = http.createServer((req, res) => {
  if (req.method !== 'POST') {
    res.writeHead(405);
    return res.end('Method Not Allowed');
  }

  let body = '';
  req.on('data', chunk => body += chunk);
  req.on('end', () => {
    try {
      const data = JSON.parse(body);
      data.thinking = { type: 'disabled' };

      // Normalize: ALL messages must have content field (DeepSeek strict)
      if (data.messages) {
        let fixed = 0;
        data.messages = data.messages.map((msg, i) => {
          if (!msg.content && msg.content !== '') {
            console.error('MSG', i, 'role=' + msg.role, 'has_tool_calls=' + !!msg.tool_calls, 'tc_len=' + (msg.tool_calls ? msg.tool_calls.length : 0));
            fixed++;
            if (msg.role === 'tool') {
              return { ...msg, content: '' };
            }
            if (msg.role === 'assistant' && msg.tool_calls) {
              return { ...msg, content: null };
            }
            return { ...msg, content: '' };
          }
          return msg;
        });
        if (fixed > 0) console.error('FIXED', fixed, 'messages missing content');
      }

      const newBody = JSON.stringify(data);

      const parsed = url.parse(req.url);
      const options = {
        hostname: DEEPSEEK_HOST,
        path: parsed.path,
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Content-Length': Buffer.byteLength(newBody),
          'Authorization': req.headers['authorization'] || '',
          'Accept': req.headers['accept'] || '',
        }
      };

      const proxyReq = (DEEPSEEK_PROTO === 'http' ? http : https).request(options, proxyRes => {
        res.writeHead(proxyRes.statusCode, proxyRes.headers);
        proxyRes.pipe(res);
      });

      proxyReq.on('error', err => {
        console.error('Proxy error:', err.message);
        res.writeHead(502);
        res.end(JSON.stringify({ error: { message: 'Proxy error: ' + err.message } }));
      });

      proxyReq.write(newBody);
      proxyReq.end();
    } catch (err) {
      console.error('Parse error:', err.message);
      res.writeHead(400);
      res.end(JSON.stringify({ error: { message: 'Invalid JSON: ' + err.message } }));
    }
  });
});

server.listen(PORT, () => {
  console.log(`DeepSeek proxy listening on port ${PORT} -> ${DEEPSEEK_PROTO}://${DEEPSEEK_HOST}`);
});
