/**
 * フロントエンド全体で使う構造化エラー。
 *
 * - `status`  : HTTP ステータスコード（ネットワーク断等は 0）
 * - `display` : エラーの表示方法。省略時は status から自動決定される
 *
 * ## status → display デフォルトマッピング
 * | status     | display    | 理由                             |
 * |------------|------------|----------------------------------|
 * | 401        | 'none'     | interceptor がリダイレクト済み   |
 * | 403        | 'dialog'   | 権限エラーは目立たせる           |
 * | 5xx / 0    | 'dialog'   | サーバー障害・ネットワーク断     |
 * | その他     | 'snackbar' | 軽微なエラー（404, 409 など）    |
 */

export type ErrorDisplay = 'snackbar' | 'dialog' | 'none'

export interface AppErrorOptions {
  /** HTTP ステータスコード。ネットワークエラー等は 0 を使う */
  status: number
  /** 表示方法を強制したい場合に指定する。省略時は status から自動決定 */
  display?: ErrorDisplay
}

function resolveDisplay(status: number): ErrorDisplay {
  if (status === 401) return 'none'
  if (status === 403) return 'dialog'
  if (status >= 500 || status === 0) return 'dialog'
  return 'snackbar'
}

export class AppError extends Error {
  readonly status: number
  readonly display: ErrorDisplay

  constructor(message: string, options: AppErrorOptions) {
    super(message)
    this.name = 'AppError'
    this.status = options.status
    this.display = options.display ?? resolveDisplay(options.status)
  }

  /**
   * 任意の unknown エラーを AppError に変換する。
   * すでに AppError の場合はそのまま返す。
   * それ以外の Error はメッセージを引き継いで status=500 として扱う。
   */
  static from(err: unknown): AppError {
    if (err instanceof AppError) return err
    const message = err instanceof Error ? err.message : '予期しないエラーが発生しました'
    return new AppError(message, { status: 500 })
  }

  static is(err: unknown): err is AppError {
    return err instanceof AppError
  }
}
