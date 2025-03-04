import {
  CanActivate,
  ExecutionContext,
  UnauthorizedException,
} from '@nestjs/common';
import * as firebase from 'firebase-admin';
import * as process from 'process';
import { Reflector } from '@nestjs/core';
import { IS_PUBLIC_KEY } from '../utils/public_annotation';

export class FirebaseMiddleware implements CanActivate {
  private firebaseAdmin: any;
  private publicEndpoints: string[] = ['/health/status'];

  constructor(private reflector: Reflector) {
    this.reflector = new Reflector();
    this.firebaseAdmin = firebase.initializeApp({
      credential: firebase.credential.cert(
        JSON.parse(process.env.FIREBASE_ADMIN_KEY),
      ),
    });
  }
  async canActivate(context: ExecutionContext): Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) {
      return true;
    }

    const req = context.switchToHttp().getRequest();
    const token = req.headers.authorization;

    if (req.method == 'OPTIONS') {
      return true;
    }
    if (this.publicEndpoints.includes(req.originalUrl)) {
      return true;
    }

    if (token != undefined && token != '') {
      return this.firebaseAdmin
        .auth()
        .verifyIdToken(token.replace('Bearer ', ''))
        .then((decodedToken: any) => {
          req.user = {
            uid: decodedToken.uid,
            email: decodedToken.email,
            name: decodedToken.name,
            picture: decodedToken.picture,
          };
          return true;
        })
        .catch((error: any) => {
          console.error('Error while verifying Firebase ID token:', error);
          this._acceddDenied(req.url);
        });
    } else {
      this._acceddDenied(req.url);
    }
  }

  _acceddDenied(url: string) {
    throw new UnauthorizedException('Unauthorized', '401 on ' + url);
  }
}
