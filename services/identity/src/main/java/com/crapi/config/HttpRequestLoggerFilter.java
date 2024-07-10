package com.crapi.config;

import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.FilterConfig;
import jakarta.servlet.ServletException;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import jakarta.servlet.http.HttpServletRequest;
import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintWriter;
import java.time.LocalDateTime;
import org.springframework.stereotype.Component;

@Component
public class HttpRequestLoggerFilter implements Filter {

  private static final String LOG_FILE_PATH = "/home/http_requests.log";

  @Override
  public void init(FilterConfig filterConfig) throws ServletException {
    // Initialization if needed
  }

  @Override
  public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
      throws IOException, ServletException {
    HttpServletRequest httpServletRequest = (HttpServletRequest) request;

    String logEntry =
        String.format(
            "Time: %s, Method: %s, URI: %s, QueryString: %s, RemoteAddr: %s%n",
            LocalDateTime.now(),
            httpServletRequest.getMethod(),
            httpServletRequest.getRequestURI(),
            httpServletRequest.getQueryString(),
            httpServletRequest.getRemoteAddr());

    // Log the request data to a file
    try (FileWriter fileWriter = new FileWriter(LOG_FILE_PATH, true);
        PrintWriter printWriter = new PrintWriter(fileWriter)) {
      printWriter.write(logEntry);
    }

    chain.doFilter(request, response);
  }

  @Override
  public void destroy() {
    // Cleanup if needed
  }
}
