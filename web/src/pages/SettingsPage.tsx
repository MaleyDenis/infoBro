import React, { useState } from 'react';

const SettingsPage: React.FC = () => {
  const [darkMode, setDarkMode] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(30);
  const [notifications, setNotifications] = useState(true);
  const [savedSources, setSavedSources] = useState([
    { id: 'reddit', name: 'Reddit', enabled: true },
    { id: 'telegram', name: 'Telegram', enabled: true },
    { id: 'rss', name: 'RSS', enabled: true },
  ]);

  const handleSourceToggle = (id: string) => {
    setSavedSources(sources =>
      sources.map(source =>
        source.id === id ? { ...source, enabled: !source.enabled } : source
      )
    );
  };

  const handleSave = () => {
    // In a real app, we would save these settings to a backend API
    alert('Settings saved successfully');
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-8">Settings</h1>
      
      <div className="bg-white rounded-lg shadow-sm overflow-hidden">
        <div className="p-6 divide-y divide-gray-200">
          {/* Appearance */}
          <div className="py-4">
            <h2 className="text-lg font-medium text-gray-900 mb-4">Appearance</h2>
            <div className="flex items-center justify-between">
              <span className="text-gray-700">Dark Mode</span>
              <button
                type="button"
                className={`${
                  darkMode ? 'bg-primary-600' : 'bg-gray-200'
                } relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2`}
                onClick={() => setDarkMode(!darkMode)}
              >
                <span
                  className={`${
                    darkMode ? 'translate-x-5' : 'translate-x-0'
                  } pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out`}
                />
              </button>
            </div>
          </div>

          {/* News Preferences */}
          <div className="py-4">
            <h2 className="text-lg font-medium text-gray-900 mb-4">News Preferences</h2>
            
            <div className="space-y-4">
              <div>
                <label htmlFor="refresh" className="block text-sm font-medium text-gray-700 mb-1">
                  Auto-refresh Interval (minutes)
                </label>
                <select
                  id="refresh"
                  value={refreshInterval}
                  onChange={(e) => setRefreshInterval(Number(e.target.value))}
                  className="input w-full rounded-md"
                >
                  <option value={0}>Never</option>
                  <option value={5}>5 minutes</option>
                  <option value={15}>15 minutes</option>
                  <option value={30}>30 minutes</option>
                  <option value={60}>1 hour</option>
                </select>
              </div>

              <div className="flex items-center justify-between">
                <span className="text-gray-700">Notifications</span>
                <button
                  type="button"
                  className={`${
                    notifications ? 'bg-primary-600' : 'bg-gray-200'
                  } relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2`}
                  onClick={() => setNotifications(!notifications)}
                >
                  <span
                    className={`${
                      notifications ? 'translate-x-5' : 'translate-x-0'
                    } pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out`}
                  />
                </button>
              </div>
            </div>
          </div>

          {/* News Sources */}
          <div className="py-4">
            <h2 className="text-lg font-medium text-gray-900 mb-4">News Sources</h2>
            
            <div className="space-y-3">
              {savedSources.map((source) => (
                <div key={source.id} className="flex items-center justify-between">
                  <span className="text-gray-700">{source.name}</span>
                  <button
                    type="button"
                    className={`${
                      source.enabled ? 'bg-primary-600' : 'bg-gray-200'
                    } relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2`}
                    onClick={() => handleSourceToggle(source.id)}
                  >
                    <span
                      className={`${
                        source.enabled ? 'translate-x-5' : 'translate-x-0'
                      } pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out`}
                    />
                  </button>
                </div>
              ))}
            </div>
          </div>
          
          {/* Save Button */}
          <div className="pt-4">
            <button
              type="button"
              onClick={handleSave}
              className="btn btn-primary w-full"
            >
              Save Settings
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;