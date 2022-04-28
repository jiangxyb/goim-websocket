<template>
  <div>
    <el-container>
      <el-header>
        <span v-show="isLogin===false" style="color: red;font-size: 30px">{{ loginTip }}</span>
        <span v-show="isLogin===true" style="color: blue;font-size: 30px">用户已登录，当前用户是：{{ loginedUser }}</span>
      </el-header>
      <el-container>
        <el-aside id="side1" width="700px">
          <!--          <el-input v-model="addr" style="margin: 10px 0px"></el-input>-->
          <!--          <el-button @click="conntosvc">建立连接</el-button>-->
          <el-input v-model="msg" style="margin: 10px 0px;" id="it" @keyup.enter="sendmsg"></el-input>
          <el-select
              v-model="value2"
              multiple
              collapse-tags
              style="margin-left: 20px;"
              placeholder="请选择发送给谁，可以多选">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value">
            </el-option>
          </el-select>
          <el-button @click="sendmsg">发送</el-button>
          <div class="clearfix" style="margin: 10px 0;text-align: center">
            <el-divider content-position="left" class="ell"></el-divider>
            <div style="text-align: center;">
              <span style="text-align: center">聊天窗口</span>
            </div>
          </div>

          <div class="bd">
            <div class="text item" v-for="(item) in msglist">
              <div class="shou">{{ item.direction }}:</div>
              <span>{{ item.dt }}</span>
              <el-divider content-position="left" class="ell"></el-divider>
            </div>
          </div>
        </el-aside>

        <el-aside class="side2" width="200px">
          <el-card class="box-card cd1">
            <div class="text item">
              <p>在线用户列表：</p>
              <p>（登录查看）</p>
            </div>
          </el-card>
          <el-card class="box-card cd2">
            <div class="text item">
              <p style="color: red"><b>{{loginedUser}}</b></p>
              <p v-for="(item) in loginlist.names">{{ item }}</p>
            </div>
          </el-card>

        </el-aside>

        <el-aside class="side2" width="200px">
          <el-card class="box-card cd1">
            <div class="text item">
              <p>已注册用户列表：</p>
            </div>
          </el-card>
          <el-card class="box-card cd2">
            <div class="text item">
              <p v-for="(item) in reglist.names">{{ item }}</p>

            </div>
          </el-card>

        </el-aside>
        <el-main>
          <div id="buttom full">
            <el-card class="box-card full">
              <p>测试账号，更多请看server根目录 user.txt</p>
              <div class="text item">
                <el-table
                    :height="200"
                    :data="tableData"
                    border
                    style="width: 100%">
                  <el-table-column
                      prop="name"
                      label="用户名"
                      width="250px">
                  </el-table-column>
                  <el-table-column
                      prop="password"
                      label="密码">
                  </el-table-column>
                </el-table>
              </div>
            </el-card>
          </div>
          <el-dialog title="登录或注册" id="dialog1" :show-close="false" :width="'100%'" :visible.sync="dialogFormVisible"
                     :close-on-click-modal="false" :modal="false" :top="'2vh'" v-show="isLogin==false">
            <el-form>
              <el-form-item label="用户名" :label-width="formLabelWidth">
                <el-input autocomplete="off" v-model="user.name"></el-input>
              </el-form-item>
              <el-form-item label="密码" :label-width="formLabelWidth">
                <el-input autocomplete="off" v-model="user.password" type="password"></el-input>
                <p style="margin-top: 20px;color: #B3C0D1">注意：长度限制：3-20</p>
              </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
              <el-button type="primary" @click="reg">注册</el-button>
              <el-button type="primary" @click="login">登 录</el-button>
            </div>
          </el-dialog>
          <el-dialog title="用户已登录" :show-close="false" :width="'100%'" :visible.sync="dialogFormVisible"
                     :close-on-click-modal="false" :modal="false" :top="'2vh'" v-show="isLogin==true">
            <div>
              <p style="color: blue;font-size: 30px">当前用户是：{{ loginedUser }}</p>
            </div>
            <div slot="footer" class="dialog-footer">
              <el-button type="primary" @click="exit">注销登录</el-button>
            </div>
          </el-dialog>
          <router-view/>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'Index',
  data() {
    return {
      loginedUser: '',
      loginTip: '请先在右下方登录再发消息!',
      isLogin: false,
      dialogFormVisible: true,
      formLabelWidth: '120px',
      user: {name: '', password: ''},
      reglist: {names: []},
      loginlist: {names: []},
      conn: {},
      msglist: [],
      msg: '请输入你要发送的内容',
      options: [{
        value: '',
        label: ''
      }],
      value2: [],
      tableData: [{
        name: 'kele',
        password: '123456'
      }, {
        name: '张三',
        password: '123'
      }, {
        name: 'namei',
        password: '678'
      }, {
        name: '王小虎',
        password: '123'
      },{
        name: '孔雀大明王',
        password: '456y78'
      }]
    };
  },
  methods: {
    pint() {
      console.log(this.value2)
    },
    reg() {
      return this.axios.post('http://localhost:9288/api/register', {
        name: this.user.name,
        password: this.user.password
      })
          .then((rsp) => {
            if (rsp.data.code == 2000) {
              this.getReglist()
            } else {
              this.$message.error(rsp.data.msg)
            }
          })
          .catch(err => {
            this.$message.error('注册失败，请重新注册');
          })
    },
    login() {
      return this.axios.post('http://localhost:9288/api/login', {
        name: this.user.name,
        password: this.user.password
      })
          .then((rsp) => {
            let obj = rsp.data
            if (obj.code == 2000) {
              console.log(obj);
              // 将token存入localStorage
              localStorage.setItem('token', JSON.stringify(obj.token))
              this.isLogin = true
              this.loginedUser = obj.name
              this.conntosvc()
              let opt = this.options
              opt.forEach(function (value, index) {
                if (value.label == obj.name) {
                  opt.splice(index, 1)
                }
              }, opt)
            } else {
              this.$message.error(rsp.data.msg)
            }
          })
          .catch(error => {
            this.$message.error('用户名或密码错误');
          });
    },
    exit() {
      localStorage.clear()
      this.isLogin = false
      console.log('开始关闭')
      this.conn.close()
      console.log('关闭结束')
      this.conn = {}
      this.msglist.push({direction: '注销', dt: '注销成功'})
      this.loginlist = {names: []}
      this.loginedUser = ''
    },
    getReglist() {
      return this.axios.get('http://localhost:9288/api/list')
          .then(rsp => {
            this.reglist.names = rsp.data.names
            let arr = []
            rsp.data.names.forEach(v => {
              let single = {value: v, label: v}
              arr.push(single)
            })
            this.options = arr
          })
    },
    getLoginlist() {
      return this.axios.get('http://localhost:9288/api/loginlist')
          .then(rsp => {
            rsp.data.names.forEach(v => {
              if (v != this.loginedUser) {
                this.loginlist.names.push(v)
              }
            })
          })
    },

    sendmsg() {
      if (this.isLogin == false) {
        this.$message.error('websocket连接建立失败,请登录再发消息!')
        return
      }
      if (this.value2.length == 0) {
        this.$message.error('请选择发送对象!')
        return
      }
      let ws = this.conn
      let wsInfo = JSON.stringify({
        type: 'normal',
        content: this.msg,
        to: this.value2
      })
      ws.send(wsInfo)
      this.msglist.push({direction: '我', dt: this.msg})
      this.msg = ''
    },
    conntosvc() {
      if (this.conn instanceof WebSocket) {
        console.log('关闭连接')
        this.conn.close()
      }
      this.conn = new WebSocket('ws://localhost:9288/api/ws')
      this.conn.onopen = () => {
        let token = JSON.parse(localStorage.getItem('token'))
        let ws = this.conn
        console.log(ws)
        ws.send(JSON.stringify({
          type: 'auth',
          content: token
        }))
      }

      this.conn.onmessage = (e) => {
        let wsInfo = JSON.parse(e.data)
        console.log('接收到消息wsInfo')
        console.log(wsInfo)

        if (wsInfo.type == 'loginlist') {

          this.loginlist.names = []
          wsInfo.content.names.forEach(v => {
            if (v != this.loginedUser) {
              console.log(v)
              this.loginlist.names.push(v)
              console.log(this.loginlist.names)
            }
          })
        }else if(wsInfo.type == 'normal') {
          let msg = {direction: wsInfo.from, dt: wsInfo.content}
          this.msglist.push(msg)
        }else if(wsInfo.type == 'ok'){
          let msg = {direction: '连接', dt: wsInfo.content}
          this.msglist.push(msg)
        }
      }
      this.conn.onclose = () => {
        console.log('close'.this.conn.readyState)
      }
      console.log('结束')
    },
    getUserInfo() {
      let token = JSON.parse(localStorage.getItem('token'))
      if (token !== '' && token !== null) {
        return this.axios({
          method: 'get',
          url: 'http://localhost:9288/api/user',
          headers: {'w-token': token},
        })
            .then(rsp => {
              if (rsp.data.code == 2000) {
                this.isLogin = true
                this.loginedUser = rsp.data.name
                this.conntosvc()
                let opt = this.options
                opt.forEach(function (value, index) {
                  if (value.label == rsp.data.name) {
                    opt.splice(index, 1)
                  }
                }, opt)
              }
            })
      }
    }
  },
  computed: {
    tip() {
      if (this.isLogin == true) {
        return '用户已登录，当前用户是：'
      } else {
        return '请先登录再发消息!'
      }
    }
  },
  created() {
    this.getReglist()
    this.getUserInfo()
  },
  mounted() {
  }

};
</script>

<style scoped>
#buttom {
  position: fixed;
  z-index: 99;
  top: 300px;
  right: 300px;
}

.el-card__body, .item {
  padding: 0;
}

p {
  margin: 0;
  line-height: 20px;
  font-size: 18px;
}

.text {
  font-size: 14px;
}

.item {
  padding: 18px 0;
}


.list, .box-card {
  margin: 0;
  height: 100%;
  line-height: 0;
  width: 100%;
}

.el-dialog__wrapper {
  position: relative;
}

#side2 {
  background-color: white;
}

.el-container {
  height: 810px;
}

.el-header, .el-footer {
  background-color: #B3C0D1;
  color: #333;
  text-align: center;
  line-height: 60px;
}

.el-aside {
  background-color: #fff;
  color: #333;
  text-align: center;
}

.el-main {
  background-color: #E9EEF3;
  color: #333;
  text-align: center;
  line-height: 30px;
}

/*左侧*/
.shou {
  text-align: left;
  margin: 0 20px;
}

.bd {
  border: black solid 2px;
  background-color: white;
  height: 590px;
  overflow: scroll;
  width: 698px;
  position: absolute;
}

.text {
  font-size: 14px;
}

.item {
  margin-bottom: 0px;
  padding: 0;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both;
}

.box-card {
  width: 100%;
  height: 500px;
}


.cd1 {
  height: 80px;
  padding: 0;
}

.cd2 {
  padding: 0;
  border: black solid 2px;
  width: 200px;
  overflow: scroll;
  height: 660px;
  position: absolute;
}

.full {
  width: 100%;
  height: 100%;
  position: relative;
  z-index: 0;
}

.ell {
  margin: 5px;
}

.el-dialog {
  margin: 0;
}

.it, .el-input__inner {
  border: #333333;
}
</style>