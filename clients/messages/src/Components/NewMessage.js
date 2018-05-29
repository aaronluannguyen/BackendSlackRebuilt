import React from "react";

export default class NewMessage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            body: ""
        }
    }

    handleSubmit(evt) {
        evt.preventDefault();
        if (this.state.body !== "") {
            let messageObj = {
                author: {
                    displayName: this.props.userInfo.displayName,
                    photoURL: this.props.userInfo.photoUrl,
                    uid: this.props.userInfo.userID
                },
                body: this.state.body
            };
            this.props.channelMessageRef.ref.push(messageObj)
                .then(() => this.setState({body: "", fbError: undefined}))
                .catch(err => this.setState({fbError: err}));
        }
    }

    render() {
        return (
            <form onSubmit={evt => this.handleSubmit(evt)}>
                {
                    this.state.fbError ?
                        <div className="alert alert-danger">
                            {this.state.fbError.message}
                        </div> :
                        undefined
                }
                <input type="text" id="text-input-bar"
                       className="form-control"
                       value={this.state.body}
                       onInput={evt => this.setState({body: evt.target.value})}
                       placeholder="Type your message..."
                />
            </form>
        );
    }
}