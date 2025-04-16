import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import LoginSuccessPage from './pages/LoginSuccessPage';
import HomePage from './pages/HomePage';
import SubscriptionsPage from './pages/SubscriptionsPage';
import StripeSuccessPage from './pages/StripeSuccessPage';
import StripeCancelPage from './pages/StripeCancelPage';
import ProtectedRoute from './components/ProtectedRoute';
import StripeDonationSuccessPage from './pages/StripeDonationSuccessPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/login-success" element={<LoginSuccessPage />} />
        <Route path="/stripe/cancel" element={<StripeCancelPage />} />
        <Route path="/stripe/success" element={<StripeSuccessPage />} />
        <Route path="/stripe/donation-success" element={<StripeDonationSuccessPage />} />
        
        {/* Protected routes */}
        <Route element={<ProtectedRoute />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/subscriptions" element={<SubscriptionsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
