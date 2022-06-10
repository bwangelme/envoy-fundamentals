from flask import Flask, request

app = Flask(__name__)

@app.route("/api")
def hello_world():
    return "<p>Hello, API Service!</p>"


@app.route("/")
def index():
    api_header = request.headers.get("API-VERSIOn")
    return "Api Service Index: %s" % api_header


def main():
    app.run(port=8001)


if __name__ == '__main__':
    main()
