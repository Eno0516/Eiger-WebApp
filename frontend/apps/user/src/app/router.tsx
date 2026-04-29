import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { UserLayout } from '@app/layouts/UserLayout'
import { RequireAuth } from '@app/guards/RequireAuth'
import { LoginPage } from '@pages/LoginPage'
import { HomePage } from '@pages/HomePage'
import { SignUpPage } from '@pages/SignUpPage'
import { EventPage } from '@pages/EventPage'

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        {/* 認証不要 */}
        <Route path="/login" element={<LoginPage />} />

        {/* 認証必須 */}
        <Route element={<RequireAuth />}>
          <Route element={<UserLayout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/sign_up" element={<SignUpPage />} />
            <Route path="/user_home/:appUUID" element={<div>ユーザーホーム（未実装）</div>} />
            <Route path="/club_page/:circleUUID" element={<div>サークルページ（未実装）</div>} />
            <Route path="/circles/:circleId/events/:eventId" element={<EventPage />} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
