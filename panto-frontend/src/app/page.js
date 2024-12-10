'use client';

import React from 'react';
const url = 'https://panto-backend-production.up.railway.app';

const Login = () => {
  const handleLogin = (provider) => {
    // Save provider to localStorage
    localStorage.setItem('provider', provider);
    

    // Redirect to appropriate OAuth endpoint
    const loginUrl =
      provider === 'github'
        ? `${url}/github`
        : `${url}/gitlab`;
    window.location.href = loginUrl;
  };

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      <h1>Login</h1>
      <button
        onClick={() => handleLogin('github')}
        style={{ margin: '10px', padding: '10px 20px' }}
      >
        Login with GitHub
      </button>
      <button
        onClick={() => handleLogin('gitlab')}
        style={{ margin: '10px', padding: '10px 20px' }}
      >
        Login with GitLab
      </button>
    </div>
  );
};

export default Login;
