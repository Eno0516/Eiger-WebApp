import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { UserLayout } from '@app/layouts/UserLayout'
import { HomePage } from '@pages/home'

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<UserLayout />}>
          <Route path="/" element={<Navigate to="/home" replace />} />
          <Route path="/home" element={<HomePage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
