import React from "react";
import Message from "./Message";

export default class ChannelMessages extends React.Component {
    render() {
        if (!this.props.channelMessageSnap) {
            return <p>Loading...</p>
        }

        let messages = [];
        this.props.channelMessageSnap.forEach(messagesSnap => {
            messages.push(<Message key={messagesSnap.key} messageSnap={messagesSnap} userInfo={this.props.userInfo}/>)
        });

        return (
            <div>
                <div>
                    {messages}
                </div>
            </div>
        )
    }
}