import { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Alert from "react-bootstrap/Alert";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import "./App.css";

function App() {
  const [youtubeChannelName, setYoutubeChannelName] = useState("");
  const [twitchChannelName, setTwitchChannelName] = useState("");
  const [youtubeError, setYoutubeError] = useState("");
  const [youtubeSuccess, setYoutubeSuccess] = useState("");
  const [twitchError, setTwitchError] = useState("");
  const [twitchSuccess, setTwitchSuccess] = useState("");

  const handleYoutubeSubmit = async (
    event: React.FormEvent<HTMLFormElement>
  ) => {
    event.preventDefault();
    setYoutubeError(""); // Clear previous error
    setYoutubeSuccess(""); // Clear previous success message
    try {
      const response = await fetch("http://localhost:8080/youtube", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ channelName: youtubeChannelName }),
      });
      if (response.ok) {
        setYoutubeSuccess("YouTube channel name sent successfully");
      } else if (response.status === 409) {
        setYoutubeError("YouTube channel is already connected");
      } else {
        setYoutubeError("Failed to send YouTube channel name");
      }
    } catch (error) {
      setYoutubeError("Error: " + error.message);
    }
  };

  const handleTwitchSubmit = async (
    event: React.FormEvent<HTMLFormElement>
  ) => {
    event.preventDefault();
    setTwitchError(""); // Clear previous error
    setTwitchSuccess(""); // Clear previous success message
    try {
      const response = await fetch("http://localhost:8080/twitch", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ channelName: twitchChannelName }),
      });
      if (response.ok) {
        setTwitchSuccess("Twitch channel name sent successfully");
      } else if (response.status === 409) {
        setTwitchError("Twitch channel is already connected");
      } else {
        setTwitchError("Failed to send Twitch channel name");
      }
    } catch (error) {
      setTwitchError("Error: " + error.message);
    }
  };

  return (
    <>
      <Container>
        <Row>
          <Col>
            <Form onSubmit={handleYoutubeSubmit}>
              <Form.Group className="mb-3">
                <Form.Label>YouTube channel</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter YouTube channel"
                  value={youtubeChannelName}
                  onChange={(e) => setYoutubeChannelName(e.target.value)}
                />
              </Form.Group>
              <Button variant="primary" type="submit">
                Submit YouTube Channel
              </Button>
            </Form>
            {youtubeError && (
              <Alert variant="danger" className="mt-3">
                {youtubeError}
              </Alert>
            )}
            {youtubeSuccess && (
              <Alert variant="success" className="mt-3">
                {youtubeSuccess}
              </Alert>
            )}
          </Col>
          <Col>
            <Form onSubmit={handleTwitchSubmit}>
              <Form.Group className="mb-3">
                <Form.Label>Twitch channel</Form.Label>
                <Form.Control
                  type="text"
                  placeholder="Enter Twitch channel"
                  value={twitchChannelName}
                  onChange={(e) => setTwitchChannelName(e.target.value)}
                />
              </Form.Group>
              <Button variant="primary" type="submit">
                Submit Twitch Channel
              </Button>
            </Form>
            {twitchError && (
              <Alert variant="danger" className="mt-3">
                {twitchError}
              </Alert>
            )}
            {twitchSuccess && (
              <Alert variant="success" className="mt-3">
                {twitchSuccess}
              </Alert>
            )}
          </Col>
        </Row>
      </Container>
    </>
  );
}

export default App;