import React from "react";
import {Link} from "react-router-dom";
import {ROUTES} from "../constants"

export default class SignInView extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            userEmail: "",
            userPassword: ""
        }
    }

    handleSignIn() {
        // Authenticate here
        this.props.history.push(ROUTES.generalChannel);
    }

    render() {
        return (
            <div className="row">
                <div className="col s4">

                </div>
                <div className="col s8">
                    <div id="form-container" className="container">
                        <div className="row">
                            <form className="col s8">
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="email" type="email" className="validate"/>
                                        <label htmlFor="email">Email</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"/>
                                        <label htmlFor="password">Password</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div>
                            <div>
                                <a className="waves-effect waves-light btn-large" onClick={() => this.handleSignIn()}>Sign In</a>
                            </div>
                            <div>
                                <p>
                                    Don't have an account yet? <Link to={ROUTES.signUp}> Sign Up </Link>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}