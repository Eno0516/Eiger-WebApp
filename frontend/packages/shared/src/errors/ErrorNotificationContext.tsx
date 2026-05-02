import { createContext, useCallback, useContext, useState, type ReactNode } from 'react'
import Alert from '@mui/material/Alert'
import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogContentText from '@mui/material/DialogContentText'
import DialogTitle from '@mui/material/DialogTitle'
import Snackbar from '@mui/material/Snackbar'
import { AppError } from './AppError'

interface ErrorNotificationContextValue {
  /**
   * AppError または任意の unknown エラーを受け取り、
   * display モードに応じてスナックバーまたはダイアログを表示する。
   *
   * TanStack Query の QueryCache / MutationCache の onError か、
   * イベントハンドラの catch 節で呼び出す。
   */
  notifyError: (err: unknown) => void
}

const ErrorNotificationContext = createContext<ErrorNotificationContextValue | null>(null)

interface NotificationState {
  open: boolean
  message: string
}

export function ErrorNotificationProvider({ children }: { children: ReactNode }) {
  const [snackbar, setSnackbar] = useState<NotificationState>({ open: false, message: '' })
  const [dialog, setDialog] = useState<NotificationState>({ open: false, message: '' })

  const notifyError = useCallback((err: unknown) => {
    const appErr = AppError.from(err)
    switch (appErr.display) {
      case 'snackbar':
        setSnackbar({ open: true, message: appErr.message })
        break
      case 'dialog':
        setDialog({ open: true, message: appErr.message })
        break
      case 'none':
        // 401 など interceptor 側で処理済み。何も表示しない。
        break
    }
  }, [])

  const closeSnackbar = useCallback(() => setSnackbar((prev) => ({ ...prev, open: false })), [])
  const closeDialog = useCallback(() => setDialog((prev) => ({ ...prev, open: false })), [])

  return (
    <ErrorNotificationContext.Provider value={{ notifyError }}>
      {children}

      {/* スナックバー — 軽微なエラー（404, 409 など） */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={closeSnackbar}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert severity="error" variant="filled" onClose={closeSnackbar}>
          {snackbar.message}
        </Alert>
      </Snackbar>

      {/* ダイアログ — 重大なエラー（403, 5xx, ネットワーク断） */}
      <Dialog open={dialog.open} onClose={closeDialog}>
        <DialogTitle>エラー</DialogTitle>
        <DialogContent>
          <DialogContentText>{dialog.message}</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDialog}>閉じる</Button>
        </DialogActions>
      </Dialog>
    </ErrorNotificationContext.Provider>
  )
}

export function useErrorNotification(): ErrorNotificationContextValue {
  const ctx = useContext(ErrorNotificationContext)
  if (!ctx) {
    throw new Error('useErrorNotification は ErrorNotificationProvider の内側で使用してください')
  }
  return ctx
}
