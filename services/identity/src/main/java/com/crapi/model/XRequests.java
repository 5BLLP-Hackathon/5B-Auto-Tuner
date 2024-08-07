package com.crapi.model;

public class XRequests {
  private String URL;
  private String UserAgent;
  private String cookies;
  private String Payload;
  private String Key;

  // Constructor
  public XRequests(String URL, String UserAgent, String cookies, String Payload, String Key) {
    this.URL = URL;
    this.UserAgent = UserAgent;
    this.cookies = cookies;
    this.Payload = Payload;
    this.Key = Key;
  }

  // Getters and Setters
  public String getURL() {
    return URL;
  }

  public void setURL(String URL) {
    this.URL = URL;
  }

  public String getUserAgent() {
    return UserAgent;
  }

  public void setUserAgent(String UserAgent) {
    this.UserAgent = UserAgent;
  }

  public String getCookies() {
    return cookies;
  }

  public void setCookies(String cookies) {
    this.cookies = cookies;
  }

  public String getPayload() {
    return Payload;
  }

  public void setPayload(String Payload) {
    this.Payload = Payload;
  }

  public String getKey() {
    return Key;
  }

  public void setKey(String Key) {
    this.Key = Key;
  }
}
