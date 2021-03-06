import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import { Top } from './Top';
import { Users } from './Users';
import { SingupForm } from './SignupForm';
import { Header } from './Header';
import { Login } from './Login';
import { Admin } from './Admin';
import { Footer } from './Footer';

export const Routes = () => (
  <Router>
    <Header />
    <Switch>
      <Route exact path='/users' component={Users} />
      <Route exact path='/signup' component={SingupForm} />
      <Route exact path='/login' component={Login} />
      <Route exact path='/admin' component={Admin} />
      <Route path='/' component={Top} />
    </Switch>
    <Footer />
  </Router>
);
