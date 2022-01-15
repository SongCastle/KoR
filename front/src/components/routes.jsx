import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import { Top } from './Top';
import { Users } from './Users';

export const Routes = () => (
  <Router>
    <Switch>
      <Route exact path='/users' component={Users} />
      <Route path='/' component={Top} />
    </Switch>
  </Router>
);
