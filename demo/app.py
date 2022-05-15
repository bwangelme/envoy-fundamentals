from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, App Service!</p>"


def main():
    app.run(port=8000)


if __name__ == '__main__':
    main()
