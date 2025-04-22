import React, { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';

// Replace static imports with lazy loaded components
const LoginPage = lazy(() => import('./pages/LoginPage'));
const LoginSuccessPage = lazy(() => import('./pages/LoginSuccessPage'));
const HomePage = lazy(() => import('./pages/HomePage'));
const SubscriptionsPage = lazy(() => import('./pages/SubscriptionsPage'));
const StripeSuccessPage = lazy(() => import('./pages/StripeSuccessPage'));
const StripeCancelPage = lazy(() => import('./pages/StripeCancelPage'));
const StripeDonationSuccessPage = lazy(() => import('./pages/StripeDonationSuccessPage'));
const AdminDashboardPage = lazy(() => import('./pages/AdminDashboardPage'));
const ClaimCodePage = lazy(() => import('./pages/ClaimCodePage'));
const OnboardingWizardPage = lazy(() => import('./pages/OnboardingWizardPage'));

// Loading component for Suspense fallback
const LoadingComponent = () => (
  <div className="flex items-center justify-center min-h-screen">
    <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
  </div>
);

function App() {
  return (
    <Router>
      <Suspense fallback={<LoadingComponent />}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/login-success" element={<LoginSuccessPage />} />
          <Route path="/stripe/success" element={<StripeSuccessPage />} />
          <Route path="/stripe/cancel" element={<StripeCancelPage />} />
          <Route path="/stripe/donation/success" element={<StripeDonationSuccessPage />} />
          <Route path="/onboarding" element={<Navigate to="/onboarding/step/0" replace />} />
          <Route path="/onboarding/step/:step" element={<OnboardingWizardPage />} />
          
          <Route element={<ProtectedRoute />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/subscriptions" element={<SubscriptionsPage />} />
            <Route path="/claim" element={<ClaimCodePage />} />
            <Route path="/claim/:code" element={<ClaimCodePage />} />
          </Route>

          <Route element={<AdminRoute />}>
            <Route path="/admin/*" element={<AdminDashboardPage />} />
          </Route>
        </Routes>
      </Suspense>
    </Router>
  );
}

export default App;
