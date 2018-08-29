import React, { Component } from 'react';
import PropTypes from 'prop-types';
import {
    Link,
    Route,
    Router,
    Switch,
} from 'react-router-dom';
import { createBrowserHistory } from 'history';
import axios from 'axios';

import {
    Col,
    Collapse,
    Container,
    DropdownItem,
    DropdownMenu,
    DropdownToggle,
    Form,
    FormGroup,
    Label,
    Input,
    Nav,
    NavItem,
    NavLink,
    Navbar,
    NavbarBrand,
    NavbarToggler,
    Row,
    UncontrolledDropdown,
} from 'reactstrap';
import 'bootstrap/dist/css/bootstrap.min.css';


const history = createBrowserHistory();


export default class App extends Component {
    render() {
        return (
            <Router history={history}>
                <Switch>
                    <Route path="/objects" component={Layout} />
                    <Route component={NotFound} />
                </Switch>
            </Router>
        );
    }
}


const Layout = () => (
    <Container fluid>
        <Row> <Col> <TopBar /> </Col> </Row>
        <Row>
            <Col sm="2"> <Menu /> </Col>
            <Col> <Body /> </Col>
        </Row>
    </Container>
);


const NotFound = () => (<React.Fragment> Not Found </React.Fragment>);


const TopBar = () => (
    <React.Fragment>
        <Navbar color="dark" dark expand>
            <NavbarBrand href="/">reactstrap</NavbarBrand>
            <Nav className="ml-auto" navbar>
                <PostSearchForm />
            </Nav>

            <Nav className="ml-auto" navbar>
                <NavItem>
                    <NavLink href="https://github.com/hellupline/hello-nurse">GitHub</NavLink>
                </NavItem>
            </Nav>
        </Navbar>
    </React.Fragment>
);


class PostSearchForm extends Component {
    // use url as value source
    state = { search: '("mobile_suit_gundam" & "sky") | "eureka_seven" - "angel" - "touhou"' }

    onChange = event => {
        const { target: { value: search } } = event;

        this.setState({ search });

        history.push(`/objects/posts/search/${search}`);
    }

    render() {
        const { search } = this.state;

        return (
            <Form inline onSubmit={e => e.preventDefault()}>
                <Input onChange={this.onChange} value={search} type="search" placeholder="Search..." />
            </Form>
        );
    }
}


const Menu = () => (
    <React.Fragment>
        <Nav vertical>
            <NavItem>
                <NavLink href="#">Link</NavLink>
            </NavItem>
            <NavItem>
                <NavLink href="#">Link</NavLink>
            </NavItem>
            <NavItem>
                <NavLink href="#">Another Link</NavLink>
            </NavItem>
            <NavItem>
                <NavLink disabled href="#">Disabled Link</NavLink>
            </NavItem>
        </Nav>
        <hr />
        <p>Link based</p>
        <Nav vertical>
            <NavLink href="#">Link</NavLink>
            <NavLink href="#">Link</NavLink>
            <NavLink href="#">Another Link</NavLink>
            <NavLink disabled href="#">Disabled Link</NavLink>
        </Nav>
    </React.Fragment>
);


const Body = () => (
    <Switch>
        <Route path="/objects/posts" component={PostControllerPrefixEndpoint} />
    </Switch>
);


const PostControllerPrefixEndpoint = ({ match }) => (
    <Switch>
        <Route path="/objects/posts/search/:query" exact component={PostIndexEndpoint} />
        <Route path="/objects/posts" exact component={PostIndexEndpoint} />
        <Route path="/objects/posts/:id" exact component={PostItemEndpoint} />
    </Switch>
);


const PostItemEndpoint = ({ match: { params: { id } } }) => (
    <PostItemController id={id} />
);


class PostItemController extends Component {
    static propTypes = {
        id: PropTypes.string.required,
    }

    render() {
        const { id } = this.props;

        return (
            <React.Fragment> PostItemController : { id } </React.Fragment>
        );
    }
}


const PostIndexEndpoint = ({ match: { params: { query } } } ) => (
    <PostIndexController query={query} />
);


class PostIndexController extends React.Component {
    static propTypes = {
        query: PropTypes.string,
    }

    static defaultProps = {
        query: '',
    }

    state = { posts: [] }

    componentDidMount() {
        const { query } = this.props;

        this.getPosts(query);
    }

    async getPosts(query) {
        if (!query) { return; }

        axios.get('http://localhost:8080/api/v1/posts', { params: { q: query } })
            .then(response => response.data)
            .then(body => body.map(({ Post }) => Post))
            .then(posts => this.setState({ posts }))
    }

    render() {
        const { query } = this.props;
        const { posts } = this.state;
        return (
            <React.Fragment>
                Search Query: { query }
                <PostIndexItemsList posts={posts} />
            </React.Fragment>
        );
    }
}


const PostIndexItemsList = ({ posts }) => (
    <React.Fragment>
        <ul>
            {posts.map(post => <PostIndexItemsRow key={post.key} post={post} /> )}
        </ul>
    </React.Fragment>
);


const PostIndexItemsRow = ({ post }) => (
    <li style={{ display: "inline-block" }}>
        <a href="#">
            <img src={`${post.value.preview_url}`} alt={post.tags.join(" ")} title={post.tags.join(" ")} />
        </a>
    </li>

);
