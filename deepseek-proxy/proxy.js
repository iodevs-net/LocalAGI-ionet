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

      // Normalize: assistant msgs with tool_calls need content:null for DeepSeek compat
      if (data.messages) {
        data.messages = data.messages.map(msg => {
          if (msg.role === 'assistant' && msg.tool_calls && !msg.content) {
            return { ...msg, content: null };
          }
          return msg;
        });
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
