import axios from 'axios'

export const httpClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

httpClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response) {
      const message = (error.response.data as { error?: string })?.error ?? error.message
      return Promise.reject(new Error(message))
    }
    return Promise.reject(error)
  },
)
