import type { InputHTMLAttributes } from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export function Input({ label, error, id, ...props }: InputProps) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
      {label && (
        <label htmlFor={id} style={{ fontSize: '0.875rem', fontWeight: 500, color: '#374151' }}>
          {label}
        </label>
      )}
      <input
        id={id}
        {...props}
        style={{
          padding: '0.5rem 0.75rem',
          borderRadius: '0.375rem',
          border: error ? '1px solid #ef4444' : '1px solid #d1d5db',
          fontSize: '0.875rem',
          outline: 'none',
          ...props.style,
        }}
      />
      {error && <span style={{ fontSize: '0.75rem', color: '#ef4444' }}>{error}</span>}
    </div>
  )
}
