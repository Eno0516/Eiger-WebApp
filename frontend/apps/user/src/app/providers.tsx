import { QueryClient, QueryClientProvider, QueryCache, MutationCache } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import { useEffect, useRef, useState, type ReactNode } from 'react'
import { GoogleOAuthProvider } from '@react-oauth/google'
import { AuthProvider } from '@eiger/shared/hooks'
import { ErrorNotificationProvider, useErrorNotification } from '@eiger/shared/errors'

const theme = createTheme()

/**
 * QueryClient を ErrorNotificationProvider の内側で生成するための内部コンポーネント。
 *
 * ErrorNotificationProvider が先にマウントされる必要があるため、
 * Providers コンポーネントと分離している。
 *
 * QueryCache / MutationCache の onError が notifyError を呼ぶことで、
 * TanStack Query 経由のすべての API エラーが自動的にスナックバー / ダイアログに表示される。
 */
function QueryClientSetup({ children }: { children: ReactNode }) {
  const { notifyError } = useErrorNotification()

  // notifyError は useCallback で安定しているが、念のため ref で最新を保持する
  const notifyErrorRef = useRef(notifyError)
  useEffect(() => {
    notifyErrorRef.current = notifyError
  }, [notifyError])

  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { staleTime: 1000 * 60, retry: 1 },
        },
        queryCache: new QueryCache({
          onError: (err) => notifyErrorRef.current(err),
        }),
        mutationCache: new MutationCache({
          onError: (err) => notifyErrorRef.current(err),
        }),
      }),
  )

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}

export function Providers({ children }: { children: ReactNode }) {
  return (
    <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_CLIENT_ID}>
      <AuthProvider>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <ErrorNotificationProvider>
            <QueryClientSetup>{children}</QueryClientSetup>
          </ErrorNotificationProvider>
        </ThemeProvider>
      </AuthProvider>
    </GoogleOAuthProvider>
  )
}
