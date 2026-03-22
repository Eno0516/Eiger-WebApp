import { useState, type FormEvent } from 'react'
import { Button } from '@eiger/shared/components'
import { Input } from '@eiger/shared/components'
import { useCreateUser } from '../hooks/useUsers'

export function CreateUserForm() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const { mutate, isPending, error } = useCreateUser()

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    mutate(
      { name, email },
      {
        onSuccess: () => {
          setName('')
          setEmail('')
        },
      },
    )
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem', maxWidth: '24rem' }}>
      <Input id="name" label="Name" value={name} onChange={(e) => setName(e.target.value)} placeholder="Alice" required />
      <Input id="email" label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="alice@example.com" required />
      {error && <p style={{ margin: 0, fontSize: '0.875rem', color: '#ef4444' }}>{error.message}</p>}
      <Button type="submit" loading={isPending}>ユーザーを追加</Button>
    </form>
  )
}
