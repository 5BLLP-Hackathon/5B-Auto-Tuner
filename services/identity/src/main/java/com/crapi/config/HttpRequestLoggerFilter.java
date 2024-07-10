package com.crapi.config;

import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.FilterConfig;
import jakarta.servlet.ReadListener;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletInputStream;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletRequestWrapper;
import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.time.LocalDateTime;
import org.springframework.stereotype.Component;

@Component
public class HttpRequestLoggerFilter implements Filter {

  private static final String LOG_FILE_PATH = "/home/http_requests.txt";

  @Override
  public void init(FilterConfig filterConfig) throws ServletException {
    // Initialization if needed
  }

  @Override
  public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
      throws IOException, ServletException {
    HttpServletRequest httpServletRequest = (HttpServletRequest) request;
    MultiReadHttpServletRequest wrappedRequest =
        new MultiReadHttpServletRequest(httpServletRequest);

    // Extract additional information
    String userAgent = wrappedRequest.getHeader("User-Agent");
    String cookie = wrappedRequest.getHeader("Cookie");
    String contentType = wrappedRequest.getContentType();
    String authorization = wrappedRequest.getHeader("Authorization");

    // Read the request payload
    StringBuilder payload = new StringBuilder();
    try (BufferedReader reader = wrappedRequest.getReader()) {
      String line;
      while ((line = reader.readLine()) != null) {
        payload.append(line).append(System.lineSeparator());
      }
    }

    String logEntry =
        String.format(
            "Time: %s, Method: %s, URL: %s, QueryString: %s, RemoteAddr: %s, UserAgent: %s, Cookie: %s, Payload: %s, ContentType: %s, Authorization: %s%n",
            LocalDateTime.now(),
            wrappedRequest.getMethod(),
            wrappedRequest.getRequestURL(),
            wrappedRequest.getQueryString(),
            wrappedRequest.getRemoteAddr(),
            userAgent,
            cookie,
            payload.toString().trim(),
            contentType,
            authorization);

    // Log the request data to a file
    try (FileWriter fileWriter = new FileWriter(LOG_FILE_PATH, true);
        PrintWriter printWriter = new PrintWriter(fileWriter)) {
      printWriter.write(logEntry);
    }

    chain.doFilter(wrappedRequest, response);
  }

  @Override
  public void destroy() {
    // Cleanup if needed
  }

  private static class MultiReadHttpServletRequest extends HttpServletRequestWrapper {

    private byte[] body;

    public MultiReadHttpServletRequest(HttpServletRequest request) throws IOException {
      super(request);
      InputStream is = request.getInputStream();
      ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
      byte[] buffer = new byte[1024];
      int len;
      while ((len = is.read(buffer)) > -1) {
        byteArrayOutputStream.write(buffer, 0, len);
      }
      body = byteArrayOutputStream.toByteArray();
    }

    @Override
    public BufferedReader getReader() throws IOException {
      return new BufferedReader(new InputStreamReader(getInputStream()));
    }

    @Override
    public ServletInputStream getInputStream() throws IOException {
      final ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(body);
      return new ServletInputStream() {
        @Override
        public boolean isFinished() {
          return byteArrayInputStream.available() == 0;
        }

        @Override
        public boolean isReady() {
          return true;
        }

        @Override
        public void setReadListener(ReadListener readListener) {}

        @Override
        public int read() throws IOException {
          return byteArrayInputStream.read();
        }
      };
    }
  }
}
