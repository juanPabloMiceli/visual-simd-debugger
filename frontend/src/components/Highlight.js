import React, { Component } from 'react';
import hljs from 'highlight.js/lib/highlight';
import javascript from 'highlight.js/lib/languages/javascript';

hljs.registerLanguage('javascript', javascript);

class Highlight extends Component {
    constructor(props) {
        super(props);
        this.nodeRef = React.createRef();
    }

    componentDidMount() {
        this.highlight();
    }

    componentDidUpdate() {
        this.highlight();
    }

    highlight = () => {
        if (this.nodeRef) {
            const nodes = this.nodeRef.current.querySelectorAll('pre');
            nodes.forEach((node) => {
                hljs.highlightBlock(node);
            });
        }
    }

    render() {
        console.log(this.props)
        
        const { content } = this.props;
        console.log(content)
        return (
            <div ref={this.nodeRef} dangerouslySetInnerHTML={{ __html: content }} />
        );
    }
}


export default Highlight;