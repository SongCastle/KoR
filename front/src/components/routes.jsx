import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

import { Top } from './Top';
import { Users } from './Users';
import { Slick } from './Slick';

export const Routes = () => (
  <Router>
    <Switch>
      <Route exact path='/users' component={Users} />
      <Route exact path='/slick' component={Slick} />
      <Route path='/' component={Top} />
    </Switch>
  </Router>
);
