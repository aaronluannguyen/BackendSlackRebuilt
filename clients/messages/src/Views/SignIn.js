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

    componentWillMount() {
        let url = AJAX.updateFLName + window.localStorage.getItem("id");
        fetch(url, {
            method: 'GET'
        })
        .then(res => res.json())
        .then(this.props.history.push(ROUTES.generalChannel))
        .catch(err => {
            if (err) {

            }
        })
    }

    handleSignIn() {
        fetch(`${AJAX.signIn}`, {
                method: 'POST',
                body: JSON.stringify(
                    {
                        email: `${this.state.userEmail}`,
                        password: `${this.state.userPassword}`
                    }
                ),
                headers: {
                    'Content-Type': `${AJAX.jsonApplication}`
                }
            }
        ).then(
            (res) => {
            if (res.status < 300) {
                let authContent = res.headers.get("Authorization");
                window.localStorage.setItem("Authorization", authContent);
                return res.json()
            } else {
                throw res
            }
        })
        .then(resJson => {
            window.localStorage.setItem("id", resJson.id);
            window.localStorage.setItem("username", resJson.userName);
            window.localStorage.setItem("firstName", resJson.firstName);
            window.localStorage.setItem("lastName", resJson.lastName);
            this.props.history.push(ROUTES.generalChannel);
        })
        .catch(error => {
           error.text().then(errMsg => {
               this.setState({error: errMsg})
           })
        });
    }

    handleSubmit(evt) {
        evt.preventDefault();
        this.handleSubmit()
    }

    render() {
        return (
            <div className="row">
                <div className="col s12">
                    <div id="form-container" className="container">
                        <div className="row">
                            <form className="col s8" onSubmit={this.handleSubmit}>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="email" type="email" className="validate"
                                            value={this.state.userEmail}
                                            onInput={evt => this.setState({userEmail: evt.target.value})}
                                        />
                                        <label htmlFor="email">Email</label>
                                    </div>
                                </div>
                                <div className="row">
                                    <div className="input-field col s12">
                                        <input id="password" type="password" className="validate"
                                           value={this.state.userPassword}
                                           onInput={evt => this.setState({userPassword: evt.target.value})}
                                        />
                                        <label htmlFor="password">Password</label>
                                    </div>
                                </div>
                            </form>
                        </div>
                        {
                            this.state.error ?
                                <div className="alert alert-danger">
                                    {this.state.error}
                                </div> :
                                undefined
                        }
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