'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import './Dashboard.css';

const backendUrl = 'https://panto-backend-production.up.railway.app';

const Dashboard = () => {
  const [userDetails, setUserDetails] = useState(null);
  const [repositories, setRepositories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [checkedRepos, setCheckedRepos] = useState({});
  const router = useRouter();

  useEffect(() => {
    const provider = localStorage.getItem('provider');
    if (!provider) {
      router.push('/'); 
      return;
    }

    const fetchData = async () => {
      try {
        // Fetch user details
        const userDetailsRes = await fetch(`${backendUrl}/${provider}/dashboard`, {
          credentials: 'include',
        });
        if (!userDetailsRes.ok) throw new Error('Failed to fetch user details');
        const userDetailsData = await userDetailsRes.json();
        setUserDetails(userDetailsData.user);

        // Fetch repositories
        const reposRes = await fetch(`${backendUrl}/${provider}/dashboard/repo`, {
          credentials: 'include',
        });
        if (!reposRes.ok) throw new Error('Failed to fetch repositories');
        const reposData = await reposRes.json();

        // Handle GitLab vs GitHub providers
        const repos = provider === 'gitlab' ? reposData.repos : reposData.user;
        setRepositories(repos);

        // Initialize checked repos state based on backend data
        const initialCheckedRepos = repos.reduce((acc, repo) => {
          const fullName = repo.full_name || repo.name;
          acc[fullName] = repo.review || false;
          return acc;
        }, {});
        setCheckedRepos(initialCheckedRepos);

      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [router]);

  const handleLogout = async () => {
    const provider = localStorage.getItem('provider');
    try {
      const response = await fetch(`${backendUrl}/${provider}/logout`, {
        method: 'GET',
        credentials: 'include',
      });
      if (!response.ok) throw new Error('Logout failed');

      localStorage.removeItem('provider');
      router.push('/');
    } catch (err) {
      console.error('Error during logout:', err);
    }
  };

  const handleCheckboxChange = async (repo) => {
    const provider = localStorage.getItem('provider');
    const fullName = repo.full_name || repo.name;
  
    setCheckedRepos(prev => ({
      ...prev,
      [fullName]: !prev[fullName]
    }));
  
    try {
      // Prepare request body based on provider
      const requestBody = provider === 'gitlab' 
        ? { id: repo.id, name: fullName }  // GitLab needs ID
        : { repoFullName: fullName };      // GitHub uses full name
  
      const response = await fetch(`${backendUrl}/${provider}/review`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });
  
      if (!response.ok) {
        throw new Error('Failed to toggle review status');
      }
  
      const data = await response.json();
      console.log('Update successful:', data);
    } catch (error) {
      console.error('Error toggling review status:', error);
      
      // Revert checkbox state on error
      setCheckedRepos(prev => ({
        ...prev,
        [fullName]: !prev[fullName]
      }));
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="dashboard-container">
      <h1>Hi, {userDetails?.name}</h1>

      {userDetails && (
        <div className="profile-logout-container">
          <img 
            src={userDetails.avatar_url} 
            alt="User Avatar" 
            className="avatar-img" 
          />
          <div onClick={handleLogout} className="logout-link">
            Logout
          </div>
        </div>
      )}

      <h2>Repos:</h2>
      {repositories.length > 0 ? (
        <div>
          {repositories.map((repo, index) => {
      const fullName = repo.full_name || repo.name;
      return (
        <div key={repo.id || index} className="repository-block">
        <div className="repo-name-container">
          <a
            href={repo.web_url || repo.html_url}
            target="_blank"
            rel="noopener noreferrer"
            className="repo-name"
          >
        {fullName}
      </a>
      </div>
      <div className="checkbox">
        <input
          type="checkbox"
          id={`checkbox-${fullName}`}
          checked={checkedRepos[fullName] || false}
          onChange={() => handleCheckboxChange(repo)}  // Pass entire repo object
        />
        <label 
          htmlFor={`checkbox-${fullName}`} 
          className="checkbox-label"
        >
          Auto Review
          </label>
        </div>
      </div>
        );
        })}
        </div>
      ) : (
        <p>No repositories found.</p>
      )}
    </div>
  );
};

export default Dashboard;
