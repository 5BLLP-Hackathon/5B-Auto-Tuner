import fs from 'fs';

export const loggingMiddleware = (store) => (next) => (action) => {
  // Log HTTP requests here
  if (action.type === 'HTTP_REQUEST') { // Ensure you define this action type in your app
    const { method, url, headers, body } = action.payload;
    const logData = {
      method,
      url,
      userAgent: headers['user-agent'],
      cookie: headers.cookie,
      payload: body,
      contentType: headers['content-type'],
      contentLanguage: headers['content-language'],
      origin: headers.origin,
      authorization: headers.authorization,
    };
    // Write to a log file
    fs.appendFile('/home/js_log.txt', JSON.stringify(logData) + '\n', (err) => {
      if (err) {
        console.error('Error writing to log file:', err);
      }
    });
  }

  return next(action);
};