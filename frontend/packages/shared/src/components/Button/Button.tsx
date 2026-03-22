import type { ButtonHTMLAttributes } from 'react'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger'
  loading?: boolean
}

export function Button({ variant = 'primary', loading, children, disabled, ...props }: ButtonProps) {
  const styles: Record<string, React.CSSProperties> = {
    primary:   { background: '#3b82f6', color: '#fff', border: 'none' },
    secondary: { background: '#e5e7eb', color: '#374151', border: 'none' },
    danger:    { background: '#ef4444', color: '#fff', border: 'none' },
  }

  return (
    <button
      {...props}
      disabled={disabled || loading}
      style={{
        padding: '0.5rem 1rem',
        borderRadius: '0.375rem',
        cursor: disabled || loading ? 'not-allowed' : 'pointer',
        fontWeight: 500,
        fontSize: '0.875rem',
        opacity: disabled || loading ? 0.6 : 1,
        ...styles[variant],
        ...props.style,
      }}
    >
      {loading ? 'Loading...' : children}
    </button>
  )
}
