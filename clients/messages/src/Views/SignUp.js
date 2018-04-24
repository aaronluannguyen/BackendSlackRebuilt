import React from "react";
import {Link} from "react-router-dom";
import {ROUTES} from "../constants"

export default class SignUpView extends React.Component {

    render() {
        return (
            <div className="row">
                <div className="col s4">

                </div>
                <div className="col s8">
                    <div id="form-container" className="container">
                        <div className="row">
                            <form className="col s12">
                                <div className="row">
                                    <div className="input-field col s6">
                                        <input id="first_name" type="text" className="validate"/>
                                        <label htmlFor="first_name">First Name</label>
                                    </div>
                                    <div className="input-field col s6">
                                        <input id="last_name" type="text" className="validate"/>
                                        <label htmlFor="last_name">Last Name</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="email" type="email" className="validate"/>
                                        <label htmlFor="email">Email</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="username" type="text" className="validate"/>
                                        <label htmlFor="username">Username</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"/>
                                        <label htmlFor="password">Password</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"/>
                                        <label htmlFor="password">Confirm Password</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div>
                            <div>
                                <a className="waves-effect waves-light btn-large">Sign Up</a>
                            </div>
                            <div>
                                Already have an account? <Link to={ROUTES.signIn}> Sign In </Link>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}