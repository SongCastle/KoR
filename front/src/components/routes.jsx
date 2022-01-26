import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import { Top } from './Top';
import { Users } from './Users';
import { SingupForm } from './SignupForm';
import { Header } from './Header';
import { Login } from './Login';

export const Routes = () => (
  <Router>
    <Header />
    <Switch>
      <Route exact path='/users' component={Users} />
      <Route exact path='/signup' component={SingupForm} />
      <Route exact path='/login' component={Login} />
      <Route path='/' component={Top} />
    </Switch>
  </Router>
);
