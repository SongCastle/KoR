import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import { Top } from './Top';
import { Users } from './Users';
import { SingupForm } from './SignupForm';

export const Routes = () => (
  <Router>
    <Switch>
      <Route exact path='/users' component={Users} />
      <Route exact path='/signup' component={SingupForm} />
      <Route path='/' component={Top} />
    </Switch>
  </Router>
);
