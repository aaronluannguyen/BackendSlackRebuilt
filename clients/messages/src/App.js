import React, { Component } from 'react';
import {HashRouter as Router, Switch, Redirect, Route} from "react-router-dom";
import {ROUTES} from "./constants";
import './App.css';
import SignUpDetailsView from "./Views/SignUpDetails";
import SignUpView from "./Views/SignUp";
import SignInView from "./Views/SignIn";
import MainView from "./Views/Main";

class App extends Component {
  render() {
    return (
      <Router>
          <Switch>
              <Route exact path={ROUTES.signIn} component={SignInView} />
              <Route path={ROUTES.signUp} component={SignUpView}/>
              <Route path={ROUTES.signUpDetails} component={SignUpDetailsView}/>
              <Route path={ROUTES.main} component={MainView}/>
              <Redirect to={ROUTES.signIn}/>
          </Switch>
      </Router>
    );
  }
}

export default App;
