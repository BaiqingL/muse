import React, { useEffect, useState } from 'react';

import Home from './Home';
import Iterate from './Iterate';

const App = () => {
  const [stage, setStage] = useState<'home' | 'iterate'>();

  useEffect(() => {
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
      const tab = tabs[0];
      console.log(tab);

      const url = new URL(tab.url!);

      if (url.hostname === 'localhost') {
        setStage('iterate');
      } else {
        setStage('home');
      }

      console.log(url);
    });
  }, []);

  return (
    <div>
      {stage === 'home' && <Home />}
      {stage === 'iterate' && <Iterate />}
    </div>
  );
};

export default App;
