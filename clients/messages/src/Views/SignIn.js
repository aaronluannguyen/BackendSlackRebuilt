import React from "react";
import {Link} from "react-router-dom";
import {ROUTES} from "../constants"
import {AJAX} from "../constants";

export default class SignInView extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            userEmail: "",
            userPassword: ""
        }
    }

    handleSignIn(evt) {
        evt.preventDefault();

        fetch(
            AJAX.signIn,
            {
                method: "POST",
                headers: {
                    "Accept": AJAX.jsonApplication,
                    "Content-Type": AJAX.jsonApplication,
                },
                body: JSON.stringify({
                    email: this.state.email,
                    password: this.state.password,
                })
            }
        ).then(res => {

        })
        .then(
            (result) => {

            },
            (error) => {

            }
        )

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
                            <form className="col s8" onSubmit={(evt) => this.handleSignIn(evt)}>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="email" type="email" className="validate"
                                            value={this.state.email}
                                            onInput={evt => this.setState({email: evt.target.value})}
                                        />
                                        <label htmlFor="email">Email</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"
                                           value={this.state.password}
                                           onInput={evt => this.setState({password: evt.target.value})}
                                        />
                                        <label htmlFor="password">Password</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div>
                            <a className="waves-effect waves-light btn-large" onClick={() => this.handleSignIn()}>Sign In</a>
                            Don't have an account yet? <Link to={ROUTES.signUp}> Sign Up </Link>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}