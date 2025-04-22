import React, { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import AdminRoute from './components/AdminRoute';
import LoadingSpinner from './components/LoadingSpinner';

// Replace static imports with lazy loaded components
const LoginPage = lazy(() => import('./pages/LoginPage'));
const LoginSuccessPage = lazy(() => import('./pages/LoginSuccessPage'));
const HomePage = lazy(() => import('./pages/HomePage'));
const SubscriptionsPage = lazy(() => import('./pages/SubscriptionsPage'));
const StripeSuccessPage = lazy(() => import('./pages/StripeSuccessPage'));
const StripeCancelPage = lazy(() => import('./pages/StripeCancelPage'));
const AdminDashboardPage = lazy(() => import('./pages/AdminDashboardPage'));
const ClaimCodePage = lazy(() => import('./pages/ClaimCodePage'));
const OnboardingWizardPage = lazy(() => import('./pages/OnboardingWizardPage'));

function App() {
  return (
    <AuthProvider>
      <Router>
        <Suspense fallback={<LoadingSpinner />}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/login-success" element={<LoginSuccessPage />} />
            <Route path="/subscription-success" element={<StripeSuccessPage type="Subscription"/>} />
            <Route path="/subscription-cancel" element={<StripeCancelPage type="Subscription" />} />
            <Route path="/donation-success" element={<StripeSuccessPage type="Donation" />} />
            <Route path="/donation-cancel" element={<StripeCancelPage type="Donation" />} />
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
    </AuthProvider>
  );
}

export default App;
