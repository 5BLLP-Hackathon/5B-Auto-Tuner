/*
 *
 * Licensed under the Apache License, Version 2.0 (the “License”);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { invalidSessionAction } from "../actions/userActions";

export const authInterceptor =
  ({ dispatch }) =>
  (next) =>
  (action) => {
    console.log(action);
    if (action.payload?.status === 401) {
      console.log("Logging out");
      return dispatch(invalidSessionAction());
    }
    if (!action.error) {
      return next(action);
    }
    return next(action);
  };

export const loggingMiddleware = (store) => (next) => (action) => {
  // Log HTTP requests here
  if (action.type === 'HTTP_REQUEST') { // You need to define this action type in your app
    const { method, url, headers, body } = action.payload;
    console.log('HTTP Request:', {
      method,
      url,
      userAgent: headers['user-agent'],
      cookie: headers.cookie,
      payload: body,
      contentType: headers['content-type'],
      contentLanguage: headers['content-language'],
      origin: headers.origin,
      authorization: headers.authorization,
    });
  }

  return next(action);
};