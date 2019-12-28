#--coding:utf-8--
from django.shortcuts import render
from django.http import HttpResponseRedirect
import hashlib
import pymysql
import time
import os
from models import LostGoods
from models import PiUser
from rest_framework.response import Response
from rest_framework.views import APIView
from serializers import *
import math

timeStamp = time.time()
timeArray = time.localtime(timeStamp)
cur_time = time.strftime("%Y-%m-%d %H:%M:%S", timeArray)

headers = {
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    "Access-Control-Allow-Origin": "*"
}


class User(APIView):
    #跨域请求必须带有，否则无法解析
    def options(self, request, format=None):
        return Response(headers=headers)
    def post(self,request,format=None):
        if request.GET['type']=='login':
            data = request.data
            # print(data)
            username = data['username'].replace(u' ',u'')
            password = data['password'].replace(u' ',u'')
            m = hashlib.md5()
            m.update(password)
            password = m.hexdigest()
            res = PiUser.objects.filter(username=username)
            # res = LostGoodsSerializers(res, many=True).data
            dic={}
            if len(res) == 0:
                dic={
                    'username':username,
                    'msg':u'用户名不存在',
                    'flag':0
                }
            else:
                userid=res[0].id
                res = res[0].password
                if res == password:
                    dic={
                        'username':username,
                        'userid':userid,
                        'msg':u'登录成功',
                        'flag':1
                    }
                else:
                    dic={
                        'username':username,
                        'userid':userid,
                        'msg':u'密码错误',
                        'flag':0
                    }
            return Response(dic,headers=headers)

        if request.GET['type']=='register':
            data = request.data
            username = data['username'].replace(u' ', u'')
            password = data['password'].replace(u' ', u'')
            repassword = data['repassword'].replace(u' ', u'')
            m = hashlib.md5()
            m.update(password)
            # 获取sign
            dic={}
            # print data
            if password != repassword:
                dic={
                    'msg':u'密码不存在，注册失败！',
                    'flag':0
                }

            elif PiUser.objects.filter(username=username):
                dic={
                    'msg':u'用户名已存在，注册失败！',
                    'flag':0
                }

            elif password == repassword:
                password = m.hexdigest()
                # sql_write("insert into pi_user(username,password,create_time) values('%s','%s','%s') " % (username, password, cur_time))
                PiUser.objects.create(username=username, password=password, create_time=cur_time)
                dic={
                    'msg':u'注册成功！',
                    'flag':1
                }
            return Response(dic,headers=headers)




class Lost(APIView):
    def options(self, request, format=None):
        return Response(headers=headers)

    def get(self, request, format=None):
    
        
        dic={}
        total_page=0
        if request.GET['type']=='all':
            serialize_object=LostGoods.objects.filter(status=1)
            dic={
                'result': LostGoodsSerializers(serialize_object, many=True).data,
            }
        elif request.GET['type']=='index':
            page=int(request.GET.get('page','1'))
            serialize_object=LostGoods.objects.filter(status=1)[(page-1)*8:page*8]
            total_page=math.ceil(LostGoods.objects.filter(status=1).count()/8.0)
            total_count=LostGoods.objects.filter(status=1).count()
            # print('***************')
            # print(math.ceil(LostGoods.objects.filter(status=1).count()/8.0))
            dic = {
                'result': LostGoodsSerializers(serialize_object, many=True).data,
                'total_page': total_page,
                'total_count':total_count
            }
        elif request.GET['type']=='detail':
            goods_id=request.GET.get('goods_id')
            # print(goods_id)
            serialize_object = LostGoods.objects.get(id=goods_id,status=1)
            serialize_object.view_count += 1
            serialize_object.save()
            dic={
                'result':LostGoodsSerializers(serialize_object).data,
            }
        elif request.GET['type']=='edit':
            goods_id=request.GET['goods_id']
            # username = data['username']
            # result = sql_select("select * from lost_goods WHERE id='%s'" % goods_id)
            result = LostGoods.objects.get(id=goods_id)
            if result.is_return == 1:
                dic={
                    'flag':0,
                    'msg':u'失物已归还禁止编辑！'
                }
            else:
                dic={
                    'flag':1,
                    'result':LostGoodsSerializers(result).data
                }
        elif request.GET['type']=='manage':
            userid=int(request.GET.get('userid'))
            result = LostGoods.objects.filter(user_id=userid,status=1)
            dic={
                'result':LostGoodsSerializers(result, many=True).data
            }
        return Response(dic,headers=headers)

    def post(self, request, format=None):

        if request.GET['type'] == 'detail':
            data = request.data
            do_return = int(data['return_goods'])
            goods_id = request.data['goods_id']  # post里面的数据不通过url过来，所以这里用request.data取数据
            # print(goods_id)
            userid = request.data['userid']
            # result = sql_select("select * from lost_goods WHERE id='%s'" % goods_id)
            result = LostGoods.objects.get(id=goods_id)

            if do_return == 1:
                # sql_write("update lost_goods set is_return=1 WHERE id='%s'"%goods_id)
                result.is_return = 1
                result.save()
                dic = {
                    'flag': 1,
                    'msg': u'归还成功',
                    'suc': 1
                }
            elif userid == int(result.user_id) and int(result.is_return) == 0:
                dic = {
                    'msg': u'失物可编辑',
                    'is_edit': 1
                }
            else:
                dic = {
                    'msg': u'用户无权编辑',
                    'is_edit': 0
                }
            return Response(dic, headers=headers)

        if request.GET['type']=='edit':
            print request.data
            f = request.FILES['file']
            data = request.data
            # print(data)
            goods_id = data['goods_id']
            goods_name = data['name']
            goods_address = data['address']
            goods_des = data['des']
            goods_phone = data['phone']
            username = data['username']
            date = time.strftime("%Y-%m-%d", time.localtime())
            second = time.strftime("%H%M%S", time.localtime())

            fileDir = "./media/file/%s/" % (date)
            if not os.path.isdir(fileDir):
                os.mkdir(fileDir, 777)

            fname = "./media/file/%s/%s.%s" % (date, second, f.name.split(".")[-1])
            fobj = open(fname, 'wb')
            for chunk in f.chunks():
                fobj.write(chunk)
            fobj.close()

            # sql_write("update lost_goods set name='%s',address='%s',pic='%s',update_time='%s',phone='%s',des='%s' WHERE id='%s' and is_return=0"\
            #           %(goods_name,goods_address,fname[1:],cur_time,goods_phone,goods_des,goods_id))
            result = LostGoods.objects.get(id=goods_id)
            if result.is_return == 0:
                result.name = goods_name
                result.address = goods_address
                result.pic = fname[1:]
                result.update_time = cur_time
                result.phone = goods_phone
                result.des = goods_des
                result.save()
                dic={
                    'flag':1,
                    'msg':u'修改成功'
                }
                # return render(request, 'edit.html', {'msg': '修改成功'})
            else:
                dic={
                    'flag':0,
                    'msg':u'失物已归还禁止编辑'
                }
                # return render(request, 'edit.html', {'msg': '失物已归还禁止编辑！'})
            return Response(dic,headers=headers)

        if request.GET['type']=='publish':
            data = request.data     #从post请求中读数据需要通过request.POST或者request.data,url里取值全都用request.GET
            print(data)
            name = data['name']
            address = data['address']
            des = data['des']
            phone = data['phone']
            user_id = data['userid']
            username = data['username']

            date = time.strftime("%Y-%m-%d", time.localtime())
            second = time.strftime("%H%M%S", time.localtime())

            fileDir = "./media/file/%s/" % (date)
            if not os.path.isdir(fileDir):
                os.mkdir(fileDir, 777)

            f = request.FILES['file']

            fname = "./media/file/%s/%s.%s" % (date, second, f.name.split(".")[-1])
            fobj = open(fname, 'wb')
            # 2MB进行切片
            for chunk in f.chunks():
                fobj.write(chunk)
            fobj.close()

            # sql_write("insert into lost_goods(name,address,pic,phone,des,\
            #        user_id,create_time,update_time,view_count,status,is_return) \
            #        values('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s')" % (
            # name, address, fname[1:], phone, des, user_id, cur_time,cur_time, 0, 1, 0))
            LostGoods.objects.create(name=name,
                                     address=address,
                                     pic=fname[1:],
                                     phone=phone,
                                     des=des,
                                     user_id=PiUser.objects.get(id=user_id),
                                     create_time=cur_time,
                                     update_time=cur_time,
                                     view_count=0,
                                     status=1,
                                     is_return=0)
            # return render(request, 'publish.html', {'msg': '发布成功', 'username': username})
            dic={
                'msg':u'发布成功',
                'username':username,
            }
            return Response(dic, headers=headers)

    def put(self,request):
        if request.GET['type'] == 'detail':
            data = request.data
            do_return = int(data['return_goods'])
            goods_id = request.data['goods_id']  # post里面的数据不通过url过来，所以这里用request.data取数据
            # print(goods_id)
            userid = request.data['userid']
            # result = sql_select("select * from lost_goods WHERE id='%s'" % goods_id)
            result = LostGoods.objects.get(id=goods_id)

            if do_return == 1:
                # sql_write("update lost_goods set is_return=1 WHERE id='%s'"%goods_id)
                result.is_return = 1
                result.save()
                dic = {
                    'flag': 1,
                    'msg': u'归还成功',
                    'suc': 1
                }
            elif userid == int(result.user_id) and int(result.is_return) == 0:
                dic = {
                    'msg': u'失物可编辑',
                    'is_edit': 1
                }
            else:
                dic = {
                    'msg': u'用户无权编辑',
                    'is_edit': 0
                }
            return Response(dic, headers=headers)

        if request.GET['type']=='edit':
            print request.data
            f = request.FILES['file']
            data = request.data
            # print(data)
            goods_id = data['goods_id']
            goods_name = data['name']
            goods_address = data['address']
            goods_des = data['des']
            goods_phone = data['phone']
            username = data['username']
            date = time.strftime("%Y-%m-%d", time.localtime())
            second = time.strftime("%H%M%S", time.localtime())

            fileDir = "./media/file/%s/" % (date)
            if not os.path.isdir(fileDir):
                os.mkdir(fileDir, 777)

            fname = "./media/file/%s/%s.%s" % (date, second, f.name.split(".")[-1])
            fobj = open(fname, 'wb')
            for chunk in f.chunks():
                fobj.write(chunk)
            fobj.close()

            # sql_write("update lost_goods set name='%s',address='%s',pic='%s',update_time='%s',phone='%s',des='%s' WHERE id='%s' and is_return=0"\
            #           %(goods_name,goods_address,fname[1:],cur_time,goods_phone,goods_des,goods_id))
            result = LostGoods.objects.get(id=goods_id)
            if result.is_return == 0:
                result.name = goods_name
                result.address = goods_address
                result.pic = fname[1:]
                result.update_time = cur_time
                result.phone = goods_phone
                result.des = goods_des
                result.save()
                dic={
                    'flag':1,
                    'msg':u'修改成功'
                }
                # return render(request, 'edit.html', {'msg': '修改成功'})
            else:
                dic={
                    'flag':0,
                    'msg':u'失物已归还禁止编辑'
                }
                # return render(request, 'edit.html', {'msg': '失物已归还禁止编辑！'})
            return Response(dic,headers=headers)

    def delete(self, request):
        dic={}
        if request.GET['type']=='manage':
            userid=request.GET['userid']
            delete_code=int(request.GET['delete_code'])
            goods_id=request.GET['goods_id']
            result=LostGoods.objects.get(id=goods_id)
            print(result.status)
            if delete_code==1:
                # print('qwe')
                result.status=0
                result.save()
                print(result.status)
                dic={
                    'msg':u'删除成功',
                    'suc':1
                }
            else:
                dic={
                    'msg':u'删除失败',
                    'suc':0
                }
        return Response(dic,headers=headers)

