from flask import Flask, render_template, request, session, jsonify
import string
import random
from PIL import Image, ImageDraw, ImageFont
import io
import base64
import uuid

app = Flask(__name__)
app.secret_key = "anam_cok_gizli_bu"

@app.route('/')
@app.route('/main')
def main_page():
    return "Main_Page"


@app.route('/register')
def register():
    session['captcha'] = random_string()
    img_str = create_captcha(session['captcha'])
    return render_template('register.html', captcha=img_str)


@app.route('/check_token', methods=['POST'])
def check_validation():
    if request.method == 'POST':
        if 'email' in request.form and 'captcha' in request.form:
            response_message = {
                'message_id': random_id()
            }
            user_email = request.form.get('email')
            user_captcha = request.form.get('captcha')
            print(f"""
                user email => {user_email},
                user captcha => {user_captcha},
                session_capctha => {session['captcha']}
            """)
            response_message['validation'] = validation_captcha({
                'user_captcha': user_captcha,
                'session_captcha': session['captcha']
            })
            session['sess_valid'] = response_message['validation']
            return jsonify(response_message), 200
        return jsonify({'message':'No data found'}), 404
    else:
        return jsonify(error="405", message="Not allowed."), 405

@app.route('/sess_check')
def get_session_valid():
    return session['sess_valid']

def random_string():
    return ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(6))

def random_id():
    digits = "123456789"
    return ''.join(random.choice(digits) for _ in range(8))

def random_uuid():
    return str(uuid.uuid4())

def create_captcha(text):
    width, height = 180, 65
    image = Image.new('RGB', (width, height), color=(73, 109, 137))
    font_path = '/usr/share/fonts/truetype/dejavu/DejaVuSerif-Bold.ttf'
    font = ImageFont.truetype(font_path, 36)
    #font = ImageFont.load_default()
    draw = ImageDraw.Draw(image)
    draw.text((10,10), text, fill=(255,255,0), font=font)
    buffered = io.BytesIO()
    image.save(buffered, format="PNG")
    img_str = base64.b64encode(buffered.getvalue()).decode("utf-8")
    #image.show()
    return img_str


def validation_captcha(data_field):
    if 'user_captcha' in data_field and 'session_captcha' in data_field:
        if data_field['user_captcha'] == data_field['session_captcha']:
            return True
    return False


if __name__ == "__main__":
    # app.run()
    app.run(host='0.0.0.0',port=5000)
