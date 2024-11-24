from flask import Flask, jsonify


app = Flask(__name__)


@app.route('/')
def home():
    return jsonify({"message": "Hello, World!"}), 200


@app.route('/healthcheck')
def healthcheck():
    return jsonify({"message": "REST API Healthy"}), 200


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
