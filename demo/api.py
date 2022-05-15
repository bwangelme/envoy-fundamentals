from flask import Flask

app = Flask(__name__)

@app.route("/api")
def hello_world():
    return "<p>Hello, API Service!</p>"


def main():
    app.run(port=8001)


if __name__ == '__main__':
    main()
